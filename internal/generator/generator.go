package generator

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/nestoca/joy/api/v1alpha1"
	"github.com/nestoca/joy/pkg/catalog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type Generator struct {
	catalogRepoGitAddr string
	catalogDir         string
	ghInstallTransport *ghinstallation.Transport
}

type Result struct {
	Release     *v1alpha1.Release `json:"release"`
	Environment *Environment      `json:"environment"`
}

type Environment struct {
	ClusterName string `json:"clusterName"`
	Namespace   string `json:"namespace"`
}

func New(catalogRepoGitAddr string, catalogDir string, githubAppId int64, githubInstallationId int64, privateKeyPath string) (*Generator, error) {
	t, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, githubAppId, githubInstallationId, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("creating github installation transport: %w", err)
	}

	generator := &Generator{
		catalogRepoGitAddr: catalogRepoGitAddr,
		catalogDir:         catalogDir,
		ghInstallTransport: t,
	}

	err = generator.init()
	if err != nil {
		return nil, fmt.Errorf("initializing generator: %w", err)
	}

	return generator, nil
}

func (r *Generator) getGitAuthenticationMethod() (*githttp.BasicAuth, error) {
	// The call to .Token will automatically renew the token if it's expired
	token, err := r.ghInstallTransport.Token(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("getting github installation token: %w", err)
	}

	return &githttp.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}, nil
}

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

func (r *Generator) Run(inputParams map[string]string) ([]*Result, error) {
	//Load Releases and relevant environment info (cluster name & namespace)
	joyCatalog, err := catalog.Load(catalog.LoadOpts{
		Dir:          r.catalogDir,
		LoadEnvs:     true,
		LoadReleases: true,
		ResolveRefs:  true,
	})
	if err != nil {
		return nil, err
	}

	var reconciledReleases []*Result
	for _, releaseGroup := range joyCatalog.Releases.Items {
		for _, release := range releaseGroup.Releases {
			if release != nil {
				log.Debug().Str("release", release.Name).Str("environment", release.Environment.Name).Msg("Processing release")
				reconciledReleases = append(reconciledReleases, &Result{
					Release: release,
					Environment: &Environment{
						ClusterName: release.Environment.Spec.Cluster,
						Namespace:   release.Environment.Spec.Namespace,
					},
				})
			}
		}
	}

	return reconciledReleases, nil
}
