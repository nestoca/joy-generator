package generator

import (
	"fmt"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	"github.com/nestoca/joy/api/v1alpha1"
	joy "github.com/nestoca/joy/pkg"
	"github.com/nestoca/joy/pkg/catalog"

	"github.com/nestoca/joy-generator/internal/github"
)

type Generator struct {
	LoadJoyContext JoyLoaderFunc
	Logger         zerolog.Logger
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

type JoyLoaderFunc func() (*JoyContext, error)

func RepoLoader(repo *github.Repo) JoyLoaderFunc {
	return func() (*JoyContext, error) {
		if err := repo.Pull(); err != nil {
			return nil, fmt.Errorf("pulling git repo: %w", err)
		}

		cfg, err := joy.LoadConfigFromCatalog(repo.Directory())
		if err != nil {
			return nil, fmt.Errorf("loading catalog config: %w", err)
		}

		cat, err := catalog.Load(repo.Directory(), cfg.KnownChartRefs())
		if err != nil {
			return nil, fmt.Errorf("loading catalog: %w", err)
		}

		return &JoyContext{Catalog: cat, Config: cfg}, nil
	}
}

// Run runs the generator and returns a slice of results. Each result contains the release, the environment where it
// will be deployed and the rendered values string.
func (generator *Generator) Run() ([]Result, error) {
	joyctx, err := generator.LoadJoyContext()
	if err != nil {
		return nil, fmt.Errorf("loading joy context: %w", err)
	}

	var reconciledReleases []Result
	for _, crossRelease := range joyctx.Catalog.Releases.Items {
		for _, release := range crossRelease.Releases {
			if release == nil {
				continue
			}

			generator.Logger.
				Debug().
				Str("release", release.Name).
				Str("environment", release.Environment.Name).
				Msg("processing release")

			chart, err := joy.ChartFromRelease(release, joyctx.Config.Charts, joyctx.Config.DefaultChartRef)
			if err != nil {
				generator.Logger.
					Error().
					Err(err).
					Str("release", release.Name).Str("environment", release.Environment.Name).
					Msgf("error getting chart for release %s", release.Name)
				continue
			}

			release.Spec.Chart.RepoUrl = chart.RepoURL
			release.Spec.Chart.Name = chart.Name
			release.Spec.Chart.Version = chart.Version

			values, err := joy.ReleaseValues(release, joyctx.Config.ValueMapping)
			if err != nil {
				generator.Logger.
					Error().
					Err(err).
					Str("release", release.Name).Str("environment", release.Environment.Name).
					Msgf("error computing values for release %s", release.Name)
				continue
			}

			renderedValues, err := yaml.Marshal(values)
			if err != nil {
				generator.Logger.
					Error().
					Err(err).
					Str("release", release.Name).Str("environment", release.Environment.Name).
					Msgf("error marshaling values for release %s", release.Name)
				continue
			}

			reconciledReleases = append(reconciledReleases, Result{
				Release:     release,
				Environment: release.Environment,
				Project:     release.Project,
				Values:      string(renderedValues),
			})
		}
	}

	return reconciledReleases, nil
}
