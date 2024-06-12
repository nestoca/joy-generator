package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"time"

	"github.com/davidmdm/x/xcontext"
	"github.com/rs/zerolog"

	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/github"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if errors.Is(err, context.Canceled) {
			return
		}
		os.Exit(1)
	}
}

func run() error {
	logger := zerolog.New(os.Stdout)

	cfg := GetConfig()

	ctx, stop := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Google.Repository != "" {
		if err := AuthenticateHelm(ctx, cfg.Google.Repository, cfg.Google.RawCredentials); err != nil {
			return fmt.Errorf("failed to authenticate to helm: %w", err)
		}
		logger.Info().
			Str("registry", cfg.Google.Repository).
			Int("credentials_length", len(cfg.Google.RawCredentials)).
			Msg("successfully authenticated to helm")
	}

	repo, err := func() (*github.Repo, error) {
		if !reflect.ValueOf(cfg.Github.App).IsZero() {
			return cfg.Github.App.NewRepo(cfg.Catalog)
		}
		return cfg.Github.User.NewRepo(cfg.Catalog)
	}()
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	logger.Info().Str("catalog_path", repo.Metadata.Path).Msg("initialized repo")

	repo = repo.WithLogger(logger)

	server := &http.Server{
		Addr: cfg.Port,
		Handler: Handler(HandlerParams{
			pluginToken: cfg.PluginToken,
			logger:      logger,
			repo:        repo,
			generator: &generator.Generator{
				CacheRoot:      cfg.CacheRoot,
				LoadJoyContext: generator.RepoLoader(repo),
				Logger:         logger,
				ChartPuller:    generator.MakeChartPuller(logger),
			},
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errChan := make(chan error, 1)

	go func() {
		logger.Info().Str("address", server.Addr).Msg("starting server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("failed to listen and serve: %w", err)
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.GracePeriod)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}

	logger.Info().Str("cause", context.Cause(ctx).Error()).Msg("server shutdown gracefully")

	return nil
}

func AuthenticateHelm(ctx context.Context, registry string, credentials []byte) error {
	login := exec.CommandContext(ctx, "helm", "registry", "login", "-u", "_json_key", "--password-stdin", registry)

	var buffer bytes.Buffer
	login.Stdout = &buffer
	login.Stderr = &buffer
	login.Stdin = bytes.NewReader(credentials)

	if err := login.Run(); err != nil {
		return fmt.Errorf("%w: %q", err, &buffer)
	}

	return nil
}
