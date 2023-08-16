package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PluginToken string `required:"true" split_words:"true"`
	RepoUrl     string `required:"true" split_words:"true"`
	CatalogDir  string `split_words:"true"`
	GithubToken string `split_words:"true"`
}

type GitHubAppConfig struct {
	Id             int64  `required:"true" split_words:"true"`
	InstallationId int64  `required:"true" split_words:"true"`
	PrivateKeyPath string `required:"true" split_words:"true"`
}

func Load() (*Config, *GitHubAppConfig, error) {
	newConfig := &Config{}
	err := envconfig.Process("joy", newConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("reading config: %w", err)
	}

	ghaConfig := &GitHubAppConfig{}
	// If the GithubToken is not set, the GitHub App configuration is required
	if newConfig.GithubToken == "" {
		err := envconfig.Process("joy_github_app", ghaConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("reading github app config: %w", err)
		}
	} else {
		ghaConfig = nil
	}

	// If the catalog directory is not set, create a temporary directory
	if newConfig.CatalogDir == "" {
		newConfig.CatalogDir, err = os.MkdirTemp("", "joy-catalog")
		if err != nil {
			panic(err)
		}
		log.Debug().Msgf("JOY_CATALOG_DIR not set, using %s", newConfig.CatalogDir)
	} else {
		log.Debug().Msgf("JOY_CATALOG_DIR set to %s", newConfig.CatalogDir)
	}

	return newConfig, ghaConfig, nil
}
