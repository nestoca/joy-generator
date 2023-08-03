package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy-generator/internal/apiserver"
	"github.com/nestoca/joy-generator/internal/config"
	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/gitrepo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	var repo *gitrepo.GitRepo
	if cfg.GithubApp != nil {
		repo, err = gitrepo.NewWithGithubApp(
			cfg.CatalogDir,
			cfg.RepoUrl,
			int64(cfg.GithubApp.AppId),
			int64(cfg.GithubApp.InstallationId),
			cfg.GithubApp.PrivateKeyPath,
		)
	} else {
		repo, err = gitrepo.NewWithGithubToken(
			cfg.CatalogDir,
			cfg.RepoUrl,
			cfg.GithubToken,
		)
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize git repo")
		os.Exit(1)
	}

	err = apiserver.New(
		generator.New(repo),
	).Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
