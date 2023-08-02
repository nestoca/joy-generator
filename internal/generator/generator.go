package generator

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
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
	gitAuthMethod      transport.AuthMethod
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

	// TODO: Automatic refresh of token
	token, err := t.Token(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("getting github installation token: %w", err)
	}

	generator := &Generator{
		catalogRepoGitAddr: catalogRepoGitAddr,
		catalogDir:         catalogDir,
		gitAuthMethod: &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
	}

	err = generator.init()
	if err != nil {
		return nil, fmt.Errorf("initializing generator: %w", err)
	}

	return generator, nil
}

func (r *Generator) init() error {
	err := os.Chdir(r.catalogDir)
	if err != nil {
		return fmt.Errorf("changing directory to %s: %w", r.catalogDir, err)
	}

	_, err = git.PlainClone(r.catalogDir, false, &git.CloneOptions{
		URL:           r.catalogRepoGitAddr,
		Auth:          r.gitAuthMethod,
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
