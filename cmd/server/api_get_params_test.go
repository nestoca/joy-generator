package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/github"
)

func TestGetParamsE2E(t *testing.T) {
	logs := &TestLogOutputs{}

	logger := zerolog.New(logs)

	user := github.User{
		Name:  os.Getenv("GH_USER"),
		Token: os.Getenv("GH_TOKEN"),
	}

	catalog := github.RepoMetadata{
		Path:           cmp.Or(os.Getenv("CATALOG_PATH"), filepath.Join(os.TempDir(), "catalog-test")),
		URL:            os.Getenv("CATALOG_URL"),
		TargetRevision: os.Getenv("CATALOG_REVISION"),
	}

	require.NoError(t, os.RemoveAll(catalog.Path))

	repo, err := user.NewRepo(catalog)
	require.NoError(t, err, "failed to create repo for user: %s", user.Name)

	repo = repo.WithLogger(logger)

	handler := Handler(HandlerParams{
		pluginToken: "test-token",
		logger:      logger,
		repo:        repo,
		generator: &generator.Generator{
			LoadJoyContext: generator.RepoLoader(repo),
			Logger:         logger,
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

	require.Greater(t, len(logs.Records), 0)
	for _, record := range logs.Records {
		require.NotEmpty(t, record["level"])
		require.NotEqual(t, "error", record["level"])
	}
}

type TestLogOutputs struct {
	Records []map[string]any
}

func (output *TestLogOutputs) Write(data []byte) (int, error) {
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		return 0, fmt.Errorf("invalid record: %w", err)
	}
	output.Records = append(output.Records, record)
	return len(data), nil
}
