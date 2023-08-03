package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Config struct {
	RepoUrl     string
	CatalogDir  string
	GithubToken string
	GithubApp   *GitHubAppConfig
}

type GitHubAppConfig struct {
	AppId          int
	InstallationId int
	PrivateKeyPath string
}

func Load() (*Config, error) {
	newConfig := &Config{}
	var ok bool
	var err error

	newConfig.RepoUrl, ok = os.LookupEnv("JOY_CATALOG_REPO_URL")
	if !ok {
		return nil, fmt.Errorf("JOY_CATALOG_REPO_URL not set")
	}

	newConfig.GithubToken, ok = os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		newConfig.GithubApp, err = loadGithubAppConfig()
		if err != nil {
			return nil, err
		}
	}

	newConfig.CatalogDir, ok = os.LookupEnv("JOY_CATALOG_DIR")
	if !ok {
		var err error
		newConfig.CatalogDir, err = os.MkdirTemp("", "joy-catalog")
		if err != nil {
			panic(err)
		}
		log.Debug().Msgf("JOY_CATALOG_DIR not set, using %s", newConfig.CatalogDir)
	} else {
		log.Debug().Msgf("JOY_CATALOG_DIR set to %s", newConfig.CatalogDir)
	}

	return newConfig, nil
}

func loadGithubAppConfig() (*GitHubAppConfig, error) {
	var err error
	newConfig := &GitHubAppConfig{}

	githubAppIdStr, ok := os.LookupEnv("GITHUB_APP_ID")
	if !ok {
		return nil, fmt.Errorf("GITHUB_APP_ID not set")
	}

	newConfig.AppId, err = strconv.Atoi(githubAppIdStr)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_ID not a number")
	}

	githubAppInstallationIdStr, ok := os.LookupEnv("GITHUB_APP_INSTALLATION_ID")
	if !ok {
		return nil, fmt.Errorf("GITHUB_APP_INSTALLATION_ID not set")
	}

	newConfig.InstallationId, err = strconv.Atoi(githubAppInstallationIdStr)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_INSTALLATION_ID not a number")
	}

	newConfig.PrivateKeyPath, ok = os.LookupEnv("GITHUB_APP_PRIVATE_KEY_PATH")
	if !ok {
		return nil, fmt.Errorf("GITHUB_APP_PRIVATE_KEY_PATH not set")
	}

	return newConfig, nil
}
