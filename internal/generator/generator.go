package generator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/nestoca/joy-generator/internal/gitrepo"
	"github.com/nestoca/joy/api/v1alpha1"
	joy "github.com/nestoca/joy/pkg"
	"github.com/nestoca/joy/pkg/catalog"
)

type Generator struct {
	loadJoyContext JoyLoaderFunc
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

func RepoLoader(repo *gitrepo.GitRepo) JoyLoaderFunc {
	return func() (*JoyContext, error) {
		if err := repo.Pull(); err != nil {
			return nil, fmt.Errorf("pulling git repo: %w", err)
		}

		cat, err := catalog.Load(catalog.LoadOpts{Dir: repo.Directory()})
		if err != nil {
			return nil, fmt.Errorf("loading catalog: %w", err)
		}

		cfg, err := joy.LoadCatalogConfig(repo.Directory())
		if err != nil {
			return nil, fmt.Errorf("loading catalog config: %w", err)
		}

		return &JoyContext{Catalog: cat, Config: cfg}, nil
	}
}

func New(load JoyLoaderFunc) *Generator {
	return &Generator{load}
}

// Run runs the generator and returns a slice of results. Each result contains the release, the environment where it
// will be deployed and the rendered values string.
func (r *Generator) Run() ([]Result, error) {
	joyctx, err := r.loadJoyContext()
	if err != nil {
		return nil, fmt.Errorf("loading joy context: %w", err)
	}

	chartURL, chartName := func() (base, name string) {
		if joyctx.Config.DefaultChart == "" {
			return
		}
		value, err := url.Parse(joyctx.Config.DefaultChart)
		if value.Scheme == "" {
			value, err = url.Parse("oci://" + joyctx.Config.DefaultChart)
		}
		if err != nil {
			return
		}
		return value.Host, strings.TrimPrefix(value.Path, "/")
	}()

	var reconciledReleases []Result
	for _, crossRelease := range joyctx.Catalog.Releases.Items {
		for _, release := range crossRelease.Releases {
			if release != nil {
				log.Debug().Str("release", release.Name).Str("environment", release.Environment.Name).Msg("processing release")

				if release.Spec.Chart.RepoUrl == "" {
					release.Spec.Chart.RepoUrl = chartURL
				}
				if release.Spec.Chart.Name == "" {
					release.Spec.Chart.Name = chartName
				}

				values, err := joy.ReleaseValues(release, release.Environment, joyctx.Config.ValueMapping)
				if err != nil {
					log.
						Error().
						Err(err).
						Str("release", release.Name).Str("environment", release.Environment.Name).
						Msgf("error computing values for release %s", release.Name)

					// we don't want to fail the whole process if rendering one release fails, so we'll just skip this one
					continue
				}

				renderedValues, err := yaml.Marshal(values)
				if err != nil {
					log.
						Error().
						Err(err).
						Str("release", release.Name).Str("environment", release.Environment.Name).
						Msgf("error marshaling values for release %s", release.Name)

					// we don't want to fail the whole process if rendering one release fails, so we'll just skip this one
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
	}

	return reconciledReleases, nil
}
