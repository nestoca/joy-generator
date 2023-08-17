package gitrepo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitRepo struct {
	// dir is the local directory where the catalog repositoryAddress should be cloned. Ex: /tmp/joy-catalog
	dir string

	// repositoryAddress is the HTTPS git address of the catalog repositoryAddress. Ex: https://github.com/my-org/joy-catalog.git
	url string

	// ref is the git ref to check out. Ex: refs/heads/master
	ref plumbing.ReferenceName

	// ghInstallTransport is the GitHub App authentication transport. It's used to generate a token that can be used to
	// authenticate git calls to the catalog repositoryAddress.
	ghAppInstallation *ghinstallation.Transport

	// githubToken is the GitHub Token used to authenticate API calls to the catalog repositoryAddress. When set, ghInstallTransport is
	// not used
	githubToken string

	// githubUser is the GitHub user used to authenticate API calls to the catalog repositoryAddress. Defaults to "x-access-token"
	githubUser string

	repository *git.Repository
}

// NewWithGithubApp creates a new GitRepo instance using GitHub App authentication
func NewWithGithubApp(url string, dir string, githubAppId int64, githubInstallationId int64, privateKeyPath string) (*GitRepo, error) {
	t, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, githubAppId, githubInstallationId, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("creating github installation transport: %w", err)
	}

	r := &GitRepo{
		dir:               dir,
		url:               url,
		ghAppInstallation: t,
	}

	if err := r.init(); err != nil {
		return nil, fmt.Errorf("initializing git repo: %w", err)
	}

	return r, nil
}

// NewWithGithubToken creates a new GitRepo instance using GitHub Token authentication
func NewWithGithubToken(url string, dir string, githubToken string, githubUser string) (*GitRepo, error) {
	r := &GitRepo{
		dir:         dir,
		url:         url,
		githubToken: githubToken,
		githubUser:  githubUser,
	}

	if err := r.init(); err != nil {
		return nil, fmt.Errorf("initializing git repo: %w", err)
	}

	return r, nil
}

func (r *GitRepo) init() error {
	auth, err := r.getCredentials()
	if err != nil {
		return fmt.Errorf("getting git credentials: %w", err)
	}

	repository, err := git.PlainClone(r.dir, false, &git.CloneOptions{
		URL:           r.url,
		ReferenceName: r.ref,
		Auth:          auth,
		Depth:         1, // Only fetch the latest commit
	})
	if err != nil {
		return fmt.Errorf("cloning git repository: %w", err)
	}

	r.repository = repository
	return nil
}

func (r *GitRepo) WithRef(ref string) *GitRepo {
	r.ref = plumbing.ReferenceName(ref)
	return r
}

func (r *GitRepo) Directory() string {
	return r.dir
}

// getCredentials returns an implementation githttp.AuthMethod that can be used to authenticate git calls
// to the catalog repositoryAddress
func (r *GitRepo) getCredentials() (*githttp.BasicAuth, error) {
	var token string
	var err error
	if r.githubToken != "" {
		token = r.githubToken
	} else if r.ghAppInstallation != nil {
		// The call to .Token will automatically renew the token if it's expired
		token, err = r.ghAppInstallation.Token(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("getting github installation token: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no github authentication method provided. Either githubToken or ghAppInstallation must be set")
	}

	user := "x-access-token"
	if r.githubUser != "" {
		user = r.githubUser
	}

	return &githttp.BasicAuth{
		Username: user,
		Password: token,
	}, nil
}

func (r *GitRepo) Pull() error {
	auth, err := r.getCredentials()
	if err != nil {
		return fmt.Errorf("getting git authentication credentials: %w", err)
	}

	w, err := r.repository.Worktree()
	if err != nil {
		return fmt.Errorf("loading git worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{
		ReferenceName: r.ref,
		Auth:          auth,
		Depth:         1, // Only fetch the latest commit
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("pulling git repo: %w", err)
	}

	return nil
}

// Status does a git fetch to ensure that the connection to the git repo is still intact. Used for pod status checks.
func (r *GitRepo) Status() error {
	auth, err := r.getCredentials()
	if err != nil {
		return fmt.Errorf("getting git authentication credentials: %w", err)
	}

	err = r.repository.Fetch(&git.FetchOptions{
		Auth:  auth,
		Depth: 1, // Only fetch the latest commit
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("fetching repo: %w", err)
	}

	return nil
}
