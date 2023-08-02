package main

import (
	"github.com/nestoca/joy-generator/internal/apiserver"
	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func main() {
	repoUrl, ok := os.LookupEnv("JOY_CATALOG_REPO_URL")
	if !ok {
		log.Error().Msg("JOY_CATALOG_REPO_URL not set")
		os.Exit(1)
	}

	githubAppIdStr, ok := os.LookupEnv("GITHUB_APP_ID")
	if !ok {
		log.Error().Msg("GITHUB_APP_ID not set")
		os.Exit(1)
	}

	githubAppId, err := strconv.Atoi(githubAppIdStr)
	if err != nil {
		log.Error().Msg("GITHUB_APP_ID not a number")
		os.Exit(1)
	}

	githubAppInstallationIdStr, ok := os.LookupEnv("GITHUB_APP_INSTALLATION_ID")
	if !ok {
		log.Error().Msg("GITHUB_APP_INSTALLATION_ID not set")
		os.Exit(1)
	}

	githubAppInstallationId, err := strconv.Atoi(githubAppInstallationIdStr)
	if err != nil {
		log.Error().Msg("GITHUB_APP_INSTALLATION_ID not a number")
		os.Exit(1)
	}

	privateKeyPath, ok := os.LookupEnv("GITHUB_APP_PRIVATE_KEY_PATH")
	if !ok {
		log.Error().Msg("GITHUB_APP_PRIVATE_KEY_PATH not set")
		os.Exit(1)
	}

	catalogDir, ok := os.LookupEnv("JOY_CATALOG_DIR")
	if !ok {
		var err error
		catalogDir, err = os.MkdirTemp("", "joy-catalog")
		if err != nil {
			panic(err)
		}
		log.Info().Msgf("JOY_CATALOG_DIR not set, using %s", catalogDir)
	} else {
		log.Info().Msgf("JOY_CATALOG_DIR set to %s", catalogDir)
	}

	gen, err := generator.New(catalogDir, repoUrl, int64(githubAppId), int64(githubAppInstallationId), privateKeyPath)
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
