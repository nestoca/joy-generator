package generator

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy-generator/internal/gitrepo"
	"github.com/nestoca/joy/api/v1alpha1"
	"github.com/nestoca/joy/pkg/catalog"
)

type Generator struct {
	repo *gitrepo.GitRepo
}

type Result struct {
	// Release holds the release's values loaded from the yaml file in the catalog
	Release *v1alpha1.Release `json:"release"`

	// Environment holds the environment info where the release will be deployed. The full spec is not loaded to minimize the payload size
	Environment *v1alpha1.Environment `json:"environment"`

	// Values is a yaml string that is the Release.spec.values rendered with any templated fields
	Values string `json:"values"`
}

func New(repo *gitrepo.GitRepo) *Generator {
	return &Generator{
		repo: repo,
	}
}

// Run runs the generator and returns a slice of results. Each result contains the release, the environment where it
// will be deployed and the rendered values string.
func (r *Generator) Run() ([]*Result, error) {
	// Make sure we have the latest catalog changes
	err := r.repo.Pull()
	if err != nil {
		return nil, fmt.Errorf("pulling git repo: %w", err)
	}

	// Load Releases and relevant environment info (cluster name & namespace)
	joyCatalog, err := catalog.Load(catalog.LoadOpts{
		Dir:          r.repo.Directory(),
		LoadEnvs:     true,
		LoadReleases: true,
		ResolveRefs:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("loading catalog: %w", err)
	}

	var reconciledReleases []*Result
	for _, crossRelease := range joyCatalog.Releases.Items {
		for _, release := range crossRelease.Releases {
			if release != nil {
				log.Debug().Str("release", release.Name).Str("environment", release.Environment.Name).Msg("processing release")

				renderedValues, err := RenderValues(release)
				if err != nil {
					log.Error().Err(err).Str("release", release.Name).Str("environment", release.Environment.Name).Msgf("error rendering values for release %s", release.Name)

					// we don't want to fail the whole process if rendering one release fails, so we'll just skip this one
					continue
				}

				reconciledReleases = append(reconciledReleases, &Result{
					Release:     release,
					Environment: release.Environment,
					Values:      renderedValues,
				})
			}
		}
	}

	return reconciledReleases, nil
}

func (r *Generator) Status() error {
	return r.repo.Status()
}
