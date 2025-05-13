package generator

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nestoca/joy/api/v1alpha1"
	"github.com/rs/zerolog"

	"github.com/nestoca/joy-generator/internal/github"
)

type ValueCache interface {
	Get(release *v1alpha1.Release) string
	Set(release *v1alpha1.Release, value string) string
	CleanupCache() error
	GetSize() int
}

type valueCache struct {
	logger       zerolog.Logger
	environments map[string]struct {
		releases map[string]string
	}
	previousRef string
	repo        github.Repository
	lock        *sync.Mutex
}

func (vc *valueCache) GetSize() int {
	counter := 0
	for _, env := range vc.environments {
		counter += len(env.releases)
	}
	return counter
}

func NewValueCache(repo github.Repository, logger zerolog.Logger) ValueCache {
	return &valueCache{
		environments: make(map[string]struct {
			releases map[string]string
		}),
		repo:   repo,
		logger: logger,
		lock:   &sync.Mutex{},
	}
}

func (vc *valueCache) Get(release *v1alpha1.Release) string {
	if release.File == nil {
		vc.logger.Error().Msg("Release file is nil, cannot get value from cache")
		return ""
	}
	vc.lock.Lock()
	defer vc.lock.Unlock()
	if envReleases, ok := vc.environments[release.Environment.Name]; ok {
		if ref, ok := envReleases.releases[release.File.Path]; ok {
			return ref
		}
	}
	return ""
}

func (vc *valueCache) Set(release *v1alpha1.Release, value string) string {
	if release.File == nil {
		vc.logger.Error().Msg("Release file is nil, cannot set value in cache")
		return value
	}
	vc.lock.Lock()
	defer vc.lock.Unlock()
	env := release.Environment.Name
	if _, ok := vc.environments[env]; !ok {
		vc.environments[env] = struct {
			releases map[string]string
		}{
			releases: make(map[string]string),
		}
	}
	vc.environments[env].releases[release.File.Path] = value
	return value
}

// CleanupCache checks for changes in the git repository and clears the cache if necessary.
func (vc *valueCache) CleanupCache() error {
	vc.lock.Lock()
	defer vc.lock.Unlock()
	currentRef, err := vc.repo.GetHeadSha()
	if err != nil {
		return fmt.Errorf("failed to get current git ref: %w", err)
	}

	if vc.previousRef == "" {
		vc.previousRef = currentRef
		vc.logger.Info().Msg("No previousRef found, skipping cache cleanup")
		return nil
	}

	if vc.previousRef == currentRef {
		vc.logger.Info().Msg("No changes detected, skipping cache sync")
		return nil
	}

	changedFiles, err := vc.repo.GetFilesChangedSince(vc.previousRef)
	if err != nil {
		return fmt.Errorf("failed to get files changed since %s: %w", vc.previousRef, err)
	}

	vc.logger.Info().Strs("changedFiles", changedFiles).Str("previousRef", vc.previousRef).Str("currentRef", currentRef).Msg("Change detected in git repository")

	for _, file := range changedFiles {
		if file == "joy.yaml" {
			vc.logger.Info().Msg("joy.yaml changed, clearing all cache")
			vc.environments = make(map[string]struct {
				releases map[string]string
			})
			vc.previousRef = currentRef
			return nil
		}
		if !strings.HasPrefix(file, "environments/") {
			continue
		}
		parts := strings.Split(file, "/")
		if len(parts) == 1 {
			continue
		}
		env := parts[1]
		if file == filepath.Join("environments", env, "env.yaml") {
			vc.logger.Info().Str("environment", env).Msg("env.yaml changed, clearing environment cache")
			delete(vc.environments, env)
		}
		delete(vc.environments[env].releases, filepath.Join(vc.repo.Directory(), file))
	}

	vc.previousRef = currentRef
	return nil
}
