package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog"
)

type User struct {
	Name  string
	Token string
}

type RepoMetadata struct {
	// Path is the local directory where the catalog repositoryAddress should be cloned. Ex: /tmp/joy-catalog
	Path string

	// URL is the HTTPS git address of the catalog repositoryAddress. Ex: https://github.com/my-org/joy-catalog.git
	URL string

	// TargetRevision is the revision we wish to check out: Ex: main
	TargetRevision string
}

type Repo struct {
	Metadata RepoMetadata

	credentials func() (*githttp.BasicAuth, error)

	repository *git.Repository

	mutex *sync.Mutex

	logger zerolog.Logger
}

// WithLogger create a shallow clone of the repo with the new logger set.
func (repo *Repo) WithLogger(logger zerolog.Logger) *Repo {
	clone := *repo
	clone.logger = logger
	return &clone
}

type App struct {
	ID             int64
	InstallationID int64
	PrivateKeyPath string
}

func (app App) NewRepo(metadata RepoMetadata) (*Repo, error) {
	transport, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, app.ID, app.InstallationID, app.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("creating github installation transport: %w", err)
	}

	getCredentials := func() (*githttp.BasicAuth, error) {
		token, err := transport.Token(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("getting github installation token: %w", err)
		}
		return &githttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}, nil
	}

	repo := &Repo{
		Metadata:    metadata,
		credentials: getCredentials,
		mutex:       &sync.Mutex{},
		logger:      zerolog.Nop(),
	}

	if err := repo.init(); err != nil {
		return nil, fmt.Errorf("initializing git repo: %w", err)
	}

	return repo, nil
}

func (user User) NewRepo(metadata RepoMetadata) (*Repo, error) {
	r := &Repo{
		Metadata: metadata,
		credentials: func() (*githttp.BasicAuth, error) {
			return &githttp.BasicAuth{
				Username: user.Name,
				Password: user.Token,
			}, nil
		},
		mutex:  &sync.Mutex{},
		logger: zerolog.Nop(),
	}

	if err := r.init(); err != nil {
		return nil, fmt.Errorf("initializing git repo: %w", err)
	}

	return r, nil
}

func (r *Repo) init() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.logger.Debug().Msg("opening git repository")

	repository, err := git.PlainOpen(r.Metadata.Path)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryNotExists) {
			return fmt.Errorf("opening git repository: %w", err)
		}
		auth, err := r.credentials()
		if err != nil {
			return fmt.Errorf("getting git credentials: %w", err)
		}
		repository, err = git.PlainClone(r.Metadata.Path, false, &git.CloneOptions{
			URL:  r.Metadata.URL,
			Auth: auth,
		})
		if err != nil {
			return fmt.Errorf("cloning git repository: %w", err)
		}
	}

	if revision := r.Metadata.TargetRevision; revision != "" {
		hash, err := repository.ResolveRevision(plumbing.Revision("refs/remotes/origin/" + revision))
		if err != nil {
			return fmt.Errorf("resolving revision %s: %w", revision, err)
		}

		worktree, err := repository.Worktree()
		if err != nil {
			return fmt.Errorf("getting worktree: %w", err)
		}

		checkoutOpts := &git.CheckoutOptions{
			Hash: *hash,
		}
		if err := worktree.Checkout(checkoutOpts); err != nil {
			return fmt.Errorf("checking out: %s: %w", revision, err)
		}
	}

	r.repository = repository
	return nil
}

func (r *Repo) Directory() string {
	return r.Metadata.Path
}

func (r *Repo) Pull() error {
	auth, err := r.credentials()
	if err != nil {
		return fmt.Errorf("getting git authentication credentials: %w", err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.logger.Debug().Msg("load git worktree")

	worktree, err := r.repository.Worktree()
	if err != nil {
		return fmt.Errorf("loading git worktree: %w", err)
	}

	r.logger.Debug().Msg("pull git repo")

	pullOpts := &git.PullOptions{
		Auth:  auth,
		Force: true,
	}

	if r.Metadata.TargetRevision != "" {
		pullOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + r.Metadata.TargetRevision)
	}

	if err := worktree.Pull(pullOpts); err == nil || errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	} else if !errors.Is(err, git.ErrNonFastForwardUpdate) {
		return err
	}

	// If non-fast-forward update, manually reset the branch to the remote HEAD
	localHead, err := r.repository.Head()
	if err != nil {
		return fmt.Errorf("getting local HEAD: %w", err)
	}

	remoteRefName := fmt.Sprintf("refs/remotes/origin/%s", strings.Split(localHead.Name().String(), "/")[2])

	remoteHeadHash, err := r.repository.ResolveRevision(plumbing.Revision(remoteRefName))
	if err != nil {
		return fmt.Errorf("resolving remote HEAD: %w", err)
	}

	if localHead.Hash() != *remoteHeadHash {
		r.logger.Debug().Msg("resetting to previous commit")
		if err := worktree.Reset(&git.ResetOptions{Commit: *remoteHeadHash, Mode: git.HardReset}); err != nil {
			return fmt.Errorf("resetting to previous commit: %w", err)
		}
	}

	return nil
}
