package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy-generator/internal/apiserver"
	"github.com/nestoca/joy-generator/internal/config"
	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/gitrepo"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cfg, ghAppCfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	var repo *gitrepo.GitRepo
	if ghAppCfg != nil {
		repo, err = gitrepo.NewWithGithubApp(
			cfg.RepoUrl,
			cfg.CatalogDir,
			ghAppCfg.Id,
			ghAppCfg.InstallationId,
			ghAppCfg.PrivateKeyPath,
		)
	} else {
		repo, err = gitrepo.NewWithGithubToken(
			cfg.RepoUrl,
			cfg.CatalogDir,
			cfg.GithubToken,
			cfg.GithubUser,
		)
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize git repo")
		os.Exit(1)
	}

	gen := generator.New(generator.RepoLoader(repo))

	server := apiserver.New(cfg.PluginToken, gen, repo)

	if err := server.Run(); err != nil {
		log.Error().Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
