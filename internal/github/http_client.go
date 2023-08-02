package github

import (
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"net/http"
)

func NewHttpClient(githubAppId int64, githubInstallationId int64, privateKeyData string) (*http.Client, error) {
	t, err := ghinstallation.New(http.DefaultTransport, githubAppId, githubInstallationId, []byte(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("creating github installation transport: %w", err)
	}

	// Create a custom http(s) client with your config
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: t,
	}

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(customClient))

	return customClient, nil
}
