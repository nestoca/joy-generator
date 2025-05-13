package generator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/nestoca/joy/api/v1alpha1"
	joy "github.com/nestoca/joy/pkg"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/nestoca/joy-generator/internal/github"
)

func TestCleanupCache(t *testing.T) {
	cachedFiles := []string{
		"environments/qa/releases/service-a.yaml",
		"environments/qa/releases/service-b.yaml",
		"environments/staging/releases/service-a.yaml",
		"environments/staging/releases/service-b.yaml",
	}
	testCases := []struct {
		name         string
		cachedFiles  []string
		changedFiles []string
		expectedSize int
	}{
		{
			name:         "full cache cleanup",
			cachedFiles:  cachedFiles,
			changedFiles: []string{"joy.yaml"},
			expectedSize: 0,
		},
		{
			name:         "env cache cleanup",
			cachedFiles:  cachedFiles,
			changedFiles: []string{"environments/staging/env.yaml"},
			expectedSize: 2,
		},
		{
			name:         "release cache cleanup",
			cachedFiles:  cachedFiles,
			changedFiles: []string{"environments/staging/releases/service-b.yaml"},
			expectedSize: 3,
		},
		{
			name:         "unrelated change",
			cachedFiles:  cachedFiles,
			changedFiles: []string{".github/CODEOWNERS"},
			expectedSize: 4,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// The hash will be mutated during the test to simulate new commit(s)
			commitHash := "d3adb33f"
			repo := github.RepositoryMock{
				DirectoryFunc: func() string {
					return "/tmp/catalog"
				},
				GetFilesChangedSinceFunc: func(sha string) ([]string, error) {
					return tt.changedFiles, nil
				},
				GetHeadShaFunc: func() (string, error) {
					return commitHash, nil
				},
			}
			cache := NewValueCache(&repo, zerolog.New(zerolog.Nop()))
			err := cache.CleanupCache()
			assert.NoError(t, err)

			assert.Equal(t, 0, cache.GetSize())

			// Simulate a new commit with diff
			commitHash = commitHash + "b"
			for _, file := range tt.cachedFiles {
				envName := strings.Split(file, "/")[1]
				cache.Set(&v1alpha1.Release{Environment: &v1alpha1.Environment{EnvironmentMetadata: v1alpha1.EnvironmentMetadata{Name: envName}}, File: &joy.YAMLFile{Path: filepath.Join(repo.Directory(), file)}}, "test")
			}

			assert.Equal(t, len(tt.cachedFiles), cache.GetSize())

			err = cache.CleanupCache()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSize, cache.GetSize())
		})
	}
}
