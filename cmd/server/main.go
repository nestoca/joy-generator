package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
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

	repo, err := func() (*github.Repo, error) {
		if !reflect.ValueOf(cfg.Github.App).IsZero() {
			return cfg.Github.App.NewRepo(cfg.Catalog)
		}
		return cfg.Github.User.NewRepo(cfg.Catalog)
	}()
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	repo = repo.WithLogger(logger)

	server := &http.Server{
		Addr: cfg.Port,
		Handler: Handler(HandlerParams{
			pluginToken: cfg.PluginToken,
			logger:      logger,
			repo:        repo,
			generator: &generator.Generator{
				Logger:         logger,
				LoadJoyContext: generator.RepoLoader(repo),
			},
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errChan := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	ctx, stop := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errChan:
		return fmt.Errorf("failed to listen and serve: %w", err)
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}

	return ctx.Err()
}
