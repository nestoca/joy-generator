package generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/yaml.v3"

	"github.com/nestoca/joy/api/v1alpha1"
	joy "github.com/nestoca/joy/pkg"
	"github.com/nestoca/joy/pkg/catalog"
	"github.com/nestoca/joy/pkg/helm"

	"github.com/nestoca/joy-generator/internal/github"
	"github.com/nestoca/joy-generator/internal/observability"
)

type Generator struct {
	CacheRoot      string
	LoadJoyContext JoyLoaderFunc
	Logger         zerolog.Logger
	ChartPuller    helm.Puller
	Concurrency    int
	ValueCache     ValueCache
	Lock           *sync.Mutex
}

type MutexMap sync.Map

func (m *MutexMap) Get(key string) *sync.Mutex {
	value, _ := (*sync.Map)(m).LoadOrStore(key, new(sync.Mutex))
	return value.(*sync.Mutex)
}

type ChartPuller struct {
	Logger  zerolog.Logger
	Mutexes *MutexMap
}

func MakeChartPuller(logger zerolog.Logger) ChartPuller {
	return ChartPuller{
		Logger:  logger,
		Mutexes: &MutexMap{},
	}
}

func (puller ChartPuller) Pull(ctx context.Context, opts helm.PullOptions) error {
	var buffer bytes.Buffer

	cli := helm.CLI{
		IO: joy.IO{
			Out: &buffer,
			Err: &buffer,
		},
	}

	url, _ := opts.Chart.ToURL()
	mutex := puller.Mutexes.Get(url.String())

	mutex.Lock()
	defer mutex.Unlock()

	if entries, err := os.ReadDir(opts.OutputDir); err == nil && len(entries) > 0 {
		// If the output directory exists and has content in it,
		// then it has been pulled by another goroutine: no need to pull the chart
		return nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat chart cache: %w", err)
	}

	if err := cli.Pull(ctx, opts); err != nil {
		return fmt.Errorf("%w: %q", err, &buffer)
	}

	puller.Logger.Info().Str("chart", url.String()).Msg("successfully pulled chart")

	return nil
}

type Result struct {
	// Release holds the release's values loaded from the yaml file in the catalog
	Release *v1alpha1.Release `json:"release"`

	// Environment holds the environment info where the release will be deployed. The full spec is not loaded to minimize the payload size
	Environment *v1alpha1.Environment `json:"environment"`

	// Project holds the project info of the release.
	Project *v1alpha1.Project `json:"project"`

	// Values is a yaml string that is the Release.spec.values rendered with any templated fields
	Values string `json:"values"`
}

type JoyContext struct {
	Catalog *catalog.Catalog
	Config  *joy.Config
}

type JoyLoaderFunc func(ctx context.Context) (*JoyContext, error)

func RepoLoader(repo github.Repository, valueCache ValueCache) JoyLoaderFunc {
	return func(ctx context.Context) (*JoyContext, error) {
		ctx, span := observability.StartTrace(ctx, "load_joy_context")
		defer span.End()

		if err := repo.Pull(ctx); err != nil {
			return nil, fmt.Errorf("pulling git repo: %w", err)
		}

		cfg, err := joy.LoadConfigFromCatalog(ctx, repo.Directory())
		if err != nil {
			return nil, fmt.Errorf("loading catalog config: %w", err)
		}

		cat, err := catalog.Load(ctx, repo.Directory(), cfg.KnownChartRefs())
		if err != nil {
			return nil, fmt.Errorf("loading catalog: %w", err)
		}

		err = valueCache.CleanupCache()
		if err != nil {
			return nil, fmt.Errorf("syncing value cache: %w", err)
		}

		return &JoyContext{Catalog: cat, Config: cfg}, nil
	}
}

// Run runs the generator and returns a slice of results. Each result contains the release, the environment where it
// will be deployed and the rendered values string.
func (generator *Generator) Run(ctx context.Context) ([]Result, error) {
	ctx, span := observability.StartTrace(ctx, "generator_run")
	defer span.End()

	_, waitSpan := observability.StartTrace(ctx, "waiting_for_lock")
	generator.Lock.Lock()
	defer generator.Lock.Unlock()
	waitSpan.End()

	joyCtx, err := generator.LoadJoyContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading joy context: %w", err)
	}

	chartCache := helm.ChartCache{
		Root:            generator.CacheRoot,
		Puller:          generator.ChartPuller,
		Refs:            joyCtx.Config.Charts,
		DefaultChartRef: joyCtx.Config.DefaultChartRef,
	}

	var releases []*v1alpha1.Release
	for _, cross := range joyCtx.Catalog.Releases.Items {
		for _, release := range cross.Releases {
			if release == nil {
				continue
			}
			releases = append(releases, release)
		}
	}

	span.SetAttributes(
		attribute.Int("release_count", len(joyCtx.Catalog.Releases.Items)),
	)

	var (
		wg                 sync.WaitGroup
		reconciledReleases = make([]Result, len(releases))
		semaphore          = make(chan struct{}, max(generator.Concurrency, 1))
	)

	for i, release := range releases {
		semaphore <- struct{}{}

		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

			ctx, span := observability.StartTrace(ctx, "release_render")
			defer span.End()
			span.SetAttributes(attribute.String("release", release.Name))

			span.SetAttributes(
				attribute.String("release", release.Name),
				attribute.String("env", release.Environment.Name),
			)

			chart, err := chartCache.GetReleaseChartFS(ctx, release)
			if err != nil {
				generator.Logger.
					Error().
					Err(err).
					Str("release", release.Name).Str("environment", release.Environment.Name).
					Msgf("error getting chart for release %s", release.Name)
				return
			}

			release.Spec.Chart.RepoUrl = chart.RepoURL
			release.Spec.Chart.Name = chart.Name
			release.Spec.Chart.Version = chart.Version

			cachedValues := generator.ValueCache.Get(release)
			cacheHit := cachedValues != ""
			span.SetAttributes(attribute.Bool("cacheHit", cacheHit))

			if !cacheHit {
				computedValues, err := joy.ComputeReleaseValues(release, chart)
				if err != nil {
					generator.Logger.
						Error().
						Err(err).
						Str("release", release.Name).Str("environment", release.Environment.Name).
						Msgf("error computing values for release %s", release.Name)
				}
				marshalledValues, err := yaml.Marshal(computedValues)
				if err != nil {
					generator.Logger.
						Error().
						Err(err).
						Str("release", release.Name).Str("environment", release.Environment.Name).
						Msgf("error marshalling values for release %s", release.Name)
					return
				}

				cachedValues = generator.ValueCache.Set(release, string(marshalledValues))
			}

			if err != nil {
				generator.Logger.
					Error().
					Err(err).
					Str("release", release.Name).Str("environment", release.Environment.Name).
					Msgf("error getting values for release %s", release.Name)
				return
			}

			generator.Logger.
				Debug().
				Str("release", release.Name).
				Str("environment", release.Environment.Name).
				Bool("cacheHit", cacheHit).
				Msg("processed release")

			reconciledReleases[i] = Result{
				Release:     release,
				Environment: release.Environment,
				Project:     release.Project,
				Values:      cachedValues,
			}
		}()
	}

	wg.Wait()

	return reconciledReleases, nil
}
