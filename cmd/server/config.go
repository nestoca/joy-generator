package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/davidmdm/conf"

	"github.com/nestoca/joy-generator/internal/github"
)

type Config struct {
	Port        string
	GracePeriod time.Duration

	PluginToken string

	Catalog github.RepoMetadata

	CacheRoot string

	Google struct {
		Repository          string
		CredentialsFilePath string
		RawCredentials      []byte
	}

	Github struct {
		User github.User
		App  github.App
	}
}

func GetConfig() Config {
	var cfg Config

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	conf.Var(conf.Environ, &cfg.CacheRoot, "CACHE_ROOT", conf.Default(filepath.Join(home, ".cache", "joy")))
	conf.Var(conf.Environ, &cfg.Port, "PORT", conf.Default(":8080"))
	conf.Var(conf.Environ, &cfg.GracePeriod, "GRACE_PERIOD", conf.Default(10*time.Second))
	conf.Var(conf.Environ, &cfg.PluginToken, "PLUGIN_TOKEN")

	conf.Var(conf.Environ, &cfg.Catalog.URL, "CATALOG_URL")
	conf.Var(conf.Environ, &cfg.Catalog.Path, "CATALOG_DIR", conf.Default(filepath.Join(os.TempDir(), "catalog")))
	conf.Var(conf.Environ, &cfg.Catalog.TargetRevision, "CATALOG_REVISION")

	conf.Var(conf.Environ, &cfg.Github.User.Token, "GH_TOKEN")
	conf.Var(conf.Environ, &cfg.Github.User.Name, "GH_USER")
	conf.Var(conf.Environ, &cfg.Github.App.ID, "GH_APP_ID")
	conf.Var(conf.Environ, &cfg.Github.App.InstallationID, "GH_APP_INSTALLATION_ID")
	conf.Var(conf.Environ, &cfg.Github.App.PrivateKeyPath, "GH_APP_PRIVATE_KEY_PATH")

	conf.Var(conf.Environ, &cfg.Google.CredentialsFilePath, "CREDENTIALS_FILE")
	conf.Var(conf.Environ, &cfg.Google.Repository, "GOOGLE_ARTIFACT_REPOSITORY")

	conf.Environ.MustParse()

	if path := cfg.Google.CredentialsFilePath; path != "" {
		creds, err := os.ReadFile(cfg.Google.CredentialsFilePath)
		if err != nil {
			panic(err)
		}
		cfg.Google.RawCredentials = creds
	}

	return cfg
}
