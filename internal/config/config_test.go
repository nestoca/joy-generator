package config_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/nestoca/joy-generator/internal/config"
)

func TestLoad(t *testing.T) {
	type TestData struct {
		Env                 map[string]string
		ExpectedConfig      *config.Config
		ExpectedGHAppConfig *config.GitHubAppConfig
		ExpectedError       error
	}

	testData := map[string]TestData{
		"with github app": {
			Env: map[string]string{
				"JOY_PLUGIN_TOKEN":                "abcdefg",
				"JOY_REPO_URL":                    "https://github.com/org/repo.git",
				"JOY_CATALOG_DIR":                 "/tmp/joy-catalog",
				"JOY_GITHUB_APP_ID":               "123456",
				"JOY_GITHUB_APP_INSTALLATION_ID":  "654321",
				"JOY_GITHUB_APP_PRIVATE_KEY_PATH": "/tmp/private-key.pem",
			},
			ExpectedError: nil,
			ExpectedConfig: &config.Config{
				PluginToken: "abcdefg",
				RepoUrl:     "https://github.com/org/repo.git",
				CatalogDir:  "/tmp/joy-catalog",
			},
			ExpectedGHAppConfig: &config.GitHubAppConfig{
				Id:             123456,
				InstallationId: 654321,
				PrivateKeyPath: "/tmp/private-key.pem",
			},
		},
		"with github token": {
			Env: map[string]string{
				"JOY_PLUGIN_TOKEN": "abcdefg",
				"JOY_REPO_URL":     "https://github.com/org/repo.git",
				"JOY_CATALOG_DIR":  "/tmp/joy-catalog",
				"JOY_GITHUB_TOKEN": "123456",
			},
			ExpectedError: nil,
			ExpectedConfig: &config.Config{
				PluginToken: "abcdefg",
				RepoUrl:     "https://github.com/org/repo.git",
				CatalogDir:  "/tmp/joy-catalog",
				GithubToken: "123456",
			},
			ExpectedGHAppConfig: nil,
		},
	}

	for name, data := range testData {
		t.Run(name, func(t *testing.T) {
			for key, value := range data.Env {
				err := os.Setenv(key, value)
				if err != nil {
					t.Fatalf("setting environment variable %s: %s", key, err)
				}
			}
			conf, ghappConfig, err := config.Load()
			if !errors.Is(err, data.ExpectedError) {
				t.Fatalf("expected error %s, got %s", data.ExpectedError, err)
			}

			if !reflect.DeepEqual(conf, data.ExpectedConfig) {
				t.Fatalf("expected config %v, got %v", data.ExpectedConfig, conf)
			}

			if !reflect.DeepEqual(ghappConfig, data.ExpectedGHAppConfig) {
				t.Fatalf("expected github app config %v, got %v", data.ExpectedGHAppConfig, ghappConfig)
			}

			os.Clearenv()
		})
	}
}
