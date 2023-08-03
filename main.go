package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy-generator/internal/apiserver"
	"github.com/nestoca/joy-generator/internal/config"
	"github.com/nestoca/joy-generator/internal/generator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	var gen *generator.Generator
	if cfg.GithubApp != nil {
		gen, err = generator.NewWithGitHubApp(
			cfg.CatalogDir,
			cfg.RepoUrl,
			int64(cfg.GithubApp.AppId),
			int64(cfg.GithubApp.InstallationId),
			cfg.GithubApp.PrivateKeyPath,
		)
	} else {
		gen, err = generator.NewWithGithubToken(
			cfg.CatalogDir,
			cfg.RepoUrl,
			cfg.GithubToken,
		)
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to create generator")
		os.Exit(1)
	}

	err = apiserver.New(gen).Run()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
