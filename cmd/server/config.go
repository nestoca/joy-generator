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

	Github struct {
		User github.User
		App  github.App
	}
}

func GetConfig() Config {
	var cfg Config

	conf.Var(conf.Environ, &cfg.Port, "PORT", conf.Default(":3000"))
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

	conf.Environ.MustParse()

	return cfg
}
