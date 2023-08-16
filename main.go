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
	cfg, ghAppCfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	var repo *gitrepo.GitRepo
	if ghAppCfg != nil {
		repo, err = gitrepo.NewWithGithubApp(
			cfg.CatalogDir,
			cfg.RepoUrl,
			ghAppCfg.Id,
			ghAppCfg.InstallationId,
			ghAppCfg.PrivateKeyPath,
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
		cfg.PluginToken,
		generator.New(repo),
	).Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
