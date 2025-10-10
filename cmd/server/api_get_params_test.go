package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/davidmdm/conf"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/github"
)

func TestGetParamsE2E(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("INTERNAL_TESTING")); !ok {
		t.Skip("skip internal test")
	}

	var (
		user     github.User
		catalog  github.RepoMetadata
		registry string
	)

	conf.Var(conf.Environ, &registry, "REGISTRY", conf.Required[string](true))
	conf.Var(conf.Environ, &catalog.Path, "CATALOG_PATH", conf.Default(filepath.Join(os.TempDir(), "catalog")))
	conf.Var(conf.Environ, &catalog.URL, "CATALOG_URL", conf.Required[string](true))
	conf.Var(conf.Environ, &catalog.TargetRevision, "CATALOG_REVISION", conf.Default("master"))
	conf.Var(conf.Environ, &user.Name, "GH_USER", conf.Required[string](true))
	conf.Var(conf.Environ, &user.Token, "GH_TOKEN", conf.Required[string](true))

	require.NoError(t, conf.Environ.Parse())

	require.NoError(t, os.RemoveAll(catalog.Path))

	repo, err := user.NewRepo(catalog)
	require.NoError(t, err, "failed to create repo for user: %s", user.Name)

	logs := &TestLogOutputs{
		Records: []map[string]any{},
		Mutex:   &sync.Mutex{},
	}
	logger := zerolog.New(logs)

	repo = repo.WithLogger(logger)

	cacheDir, err := os.MkdirTemp("", "joy-cache-*")
	require.NoError(t, err)

	t.Logf("cache dir: %s", cacheDir)

	valueCache := generator.NewValueCache(repo, logger)
	handler := Handler(HandlerParams{
		pluginToken: "test-token",
		logger:      logger,
		repo:        repo,
		generator: &generator.Generator{
			CacheRoot:      cacheDir,
			LoadJoyContext: generator.RepoLoader(repo, valueCache),
			Logger:         logger,
			ChartPuller:    generator.MakeChartPuller(logger),
			Concurrency:    4,
			ValueCache:     valueCache,
			Lock:           &sync.Mutex{},
		},
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL+"/api/v1/getparams.execute", strings.NewReader("{}"))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer test-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	require.NoError(t, err)

	require.Equal(t, 200, resp.StatusCode, body.String())

	var response generator.GetParamsResponse
	require.NoError(t, json.Unmarshal(body.Bytes(), &response))
	require.Greater(t, len(response.Output.Parameters), 0)

	for _, result := range response.Output.Parameters {
		chart := jsonUnmarshalTo[map[string]string](t, result.Release.Spec.Chart)
		require.NotEmpty(t, chart["version"])
		require.NotEmpty(t, chart["repoUrl"])
		require.NotEmpty(t, chart["name"])
		require.NotEmpty(t, result.Links)
	}

	require.Greater(t, len(logs.Records), 0)
	for _, record := range logs.Records {
		require.NotEmpty(t, record["level"])
		require.NotEqualf(t, "error", record["level"], "unexpected error log: %+v", record)
	}

	entries, err := os.ReadDir(cacheDir)
	require.NoError(t, err)
	require.Greater(t, len(entries), 0)
}

type TestLogOutputs struct {
	Records []map[string]any
	Mutex   *sync.Mutex
}

func (output *TestLogOutputs) Write(data []byte) (int, error) {
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		return 0, fmt.Errorf("invalid record: %w", err)
	}

	output.Mutex.Lock()
	defer output.Mutex.Unlock()

	output.Records = append(output.Records, record)
	return len(data), nil
}

func jsonUnmarshalTo[T any](t *testing.T, value any) T {
	t.Helper()

	data, err := json.Marshal(value)
	require.NoError(t, err)

	var result T
	require.NoError(t, json.Unmarshal(data, &result))

	return result
}
