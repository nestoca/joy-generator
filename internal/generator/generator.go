package generator

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy/api/v1alpha1"
	"github.com/nestoca/joy/pkg/catalog"
)

type Generator struct {
	// catalogRepoGitAddr is the HTTPS git address of the catalog repo. Ex: https://github.com/my-org/joy-catalog.git
	catalogRepoGitAddr string

	// catalogDir is the local directory where the catalog repo should be cloned. Ex: /tmp/joy-catalog
	catalogDir string

	// ghInstallTransport is the GitHub App authentication transport. It's used to generate a token that can be used to
	// authenticate git calls to the catalog repo.
	ghAppInstallation *ghinstallation.Transport

	// githubToken is the GitHub Token used to authenticate API calls to the catalog repo. When set, ghInstallTransport is
	// not used
	githubToken string
}

type Result struct {
	// Release holds the release's values loaded from the yaml file in the catalog
	Release *v1alpha1.Release `json:"release"`

	// Environment holds the environment info where the release will be deployed. The full spec is not loaded to minimize the payload size
	Environment *Environment `json:"environment"`

	// RenderedValues is a yaml string that is the Release.spec.values rendered with any templated fields
	RenderedValues string `json:"templatedValues"`
}

type Environment struct {
	Name        string `json:"name"`
	ClusterName string `json:"clusterName"`
	Namespace   string `json:"namespace"`
}

func NewWithGithubToken(catalogRepoGitAddr string, catalogDir string, githubToken string) (*Generator, error) {
	generator := &Generator{
		catalogRepoGitAddr: catalogRepoGitAddr,
		catalogDir:         catalogDir,
		githubToken:        githubToken,
	}

	err := generator.init()
	if err != nil {
		return nil, fmt.Errorf("initializing generator: %w", err)
	}

	return generator, nil
}

// NewWithGitHubApp creates a new Generator instance using GitHub App authentication
func NewWithGitHubApp(catalogRepoGitAddr string, catalogDir string, githubAppId int64, githubInstallationId int64, privateKeyPath string) (*Generator, error) {
	t, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, githubAppId, githubInstallationId, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("creating github installation transport: %w", err)
	}

	generator := &Generator{
		catalogRepoGitAddr: catalogRepoGitAddr,
		catalogDir:         catalogDir,
		ghAppInstallation:  t,
	}

	err = generator.init()
	if err != nil {
		return nil, fmt.Errorf("initializing generator: %w", err)
	}

	return generator, nil
}

// getGitAuthenticationMethod returns an implementation githttp.AuthMethod that can be used to authenticate git calls
// to the catalog repo
func (r *Generator) getGitAuthenticationMethod() (*githttp.BasicAuth, error) {
	var token string
	var err error
	if r.githubToken != "" {
		token = r.githubToken
	} else if r.ghAppInstallation != nil {
		// The call to .Token will automatically renew the token if it's expired
		token, err = r.ghAppInstallation.Token(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("getting github installation token: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no github authentication method provided. Either githubToken or ghAppInstallation must be set")
	}

	// The call to .Token will automatically renew the token if it's expired
	token, err = r.ghAppInstallation.Token(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("getting github installation token: %w", err)
	}

	return &githttp.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}, nil
}

// init initializes the generator by changing the working directory to the catalog repo and cloning it
// This should only be called once, before any calls to Run.
func (r *Generator) init() error {
	err := os.Chdir(r.catalogDir)
	if err != nil {
		return fmt.Errorf("changing directory to %s: %w", r.catalogDir, err)
	}

	auth, err := r.getGitAuthenticationMethod()
	if err != nil {
		return fmt.Errorf("getting git authentication credentials: %w", err)
	}

	_, err = git.PlainClone(r.catalogDir, false, &git.CloneOptions{
		URL:           r.catalogRepoGitAddr,
		Auth:          auth,
		ReferenceName: "merge-values-into-releases",
	})

	return err
}

// Run runs the generator and returns a slice of results. Each result contains the release, the environment where it
// will be deployed and the rendered values string.
func (r *Generator) Run() ([]*Result, error) {
	// Load Releases and relevant environment info (cluster name & namespace)
	joyCatalog, err := catalog.Load(catalog.LoadOpts{
		Dir:          r.catalogDir,
		LoadEnvs:     true,
		LoadReleases: true,
		ResolveRefs:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("loading catalog: %w", err)
	}

	var reconciledReleases []*Result
	for _, releaseGroup := range joyCatalog.Releases.Items {
		for _, release := range releaseGroup.Releases {
			if release != nil {
				log.Debug().Str("release", release.Name).Str("environment", release.Environment.Name).Msg("processing release")

				renderedValues, err := RenderValues(release)
				if err != nil {
					log.Error().Err(err).Str("release", release.Name).Str("environment", release.Environment.Name).Msgf("error rendering values for release %s", release.Name)

					// we don't want to fail rendering all the releases if rendering one fails, so we'll just skip this one
					continue
				}

				reconciledReleases = append(reconciledReleases, &Result{
					Release: release,
					Environment: &Environment{
						Name:        release.Environment.Name,
						ClusterName: release.Environment.Spec.Cluster,
						Namespace:   release.Environment.Spec.Namespace,
					},
					RenderedValues: renderedValues,
				})
			}
		}
	}

	return reconciledReleases, nil
}
