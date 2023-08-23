package repository

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"templ/configelements"
)

// Errors
type ErrEmptyUpstream struct{}

func (e ErrEmptyUpstream) Error() string {
	return "no upstream provided"
}

type ErrInvalidUpstream struct{}

func (e ErrInvalidUpstream) Error() string {
	return "invalid upstream: does not match expected patterns."
}

type ErrNotDirectory struct{}

func (e ErrNotDirectory) Error() string {
	return "not a directory"
}

// Abstract behaviours
// The Repository type implements the CommonRepositoryBehaviour interface.
type CommonRepositoryBehaviour interface {
	//TemplDestination returns the location on the file system the repository is intended for.
	TemplDestination() string
	//Fetch performs a git clone.
	Fetch() error
	//Origin returns the origin you initialised the repository with.
	Origin() string
}

type Repository struct {
	//The directory that the repository will be cloned to. The taxonomy is:
	//${TEMPL_DIR}/<repo type>/<repo identifiers, can be multiple directories>/
	Destination string
	//The location we clone from
	Upstream string
	RepoType string
	//Validation patterns for the upstream. We also use data from the upstream to generate the destination
	ValidUpstreams []string
}

// NewGitRepository handles figuring out which of the types of git repository we're dealing with and returns
// the appropriate concrete type.
// This constructor and the NewLocalGitRepository structure both benefit from checking os.Stat on the upstream string.
// but they have different purposes. Here, we're trying to determine the right repository type; in NewLocalGitRepository
// we're looking for errors.
func NewGitRepository(upstream string) (CommonRepositoryBehaviour, error) {

	if upstream == "" {
		return nil, ErrInvalidUpstream{}
	}

	stat, err := os.Stat(upstream)

	if err != nil {
		//It's a remote upstream, not a local directory
		if os.IsNotExist(err) {
			// Assuming it's a remote upstream if the path doesn't exist
			return NewGitHubRepository(upstream)
		}
	}

	if stat.IsDir() {
		return NewLocalGitRepository(upstream)
	}

	return nil, fmt.Errorf("unhandled scenario for upstream: %s", upstream)
}

func NewGitHubRepository(upstream string) (CommonRepositoryBehaviour, error) {
	if upstream == "" {
		return nil, ErrInvalidUpstream{}
	}

	repo := Repository{
		Upstream: upstream,
		ValidUpstreams: []string{
			// Regular expression pattern for GitHub repository extraction
			`(.*)[:/]([^/]+)/([^/.]+)(\.git)?$`,        // git@github.com:<owner>/<repo> or https://github.com/<owner>/<repo>
			`https?://(.*)/([^/]+)/([^/.]+)(\.git)?$`,  // https://github.com/<owner>/<repo>
			`git://(.*)/([^/]+)/([^/.]+)(\.git)?$`,     // git://github.com/<owner>/<repo>
			`ssh://git@(.*)/([^/]+)/([^/.]+)(\.git)?$`, // ssh://git@github.com/<owner>/<repo>
			`git@(.*):([^/]+)/([^/.]+)(\.git)?$`,       // git@github.com:<owner>/<repo>.git
		},
	}

	var username_reponame string

	for _, pattern := range repo.ValidUpstreams {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(repo.Upstream)
		if matches != nil {
			// Return the repository name
			username_reponame = matches[2] + "/" + matches[3]
		}
	}

	//Need to extract the owner and repo from above
	t := configelements.NewTemplDir()
	repo.Destination = fmt.Sprintf("%s/github/%s", t.TemplatesDir, username_reponame)

	return repo, nil
}

func NewLocalGitRepository(upstream string) (CommonRepositoryBehaviour, error) {
	if upstream == "" {
		return nil, ErrInvalidUpstream{}
	}

	stat, err := os.Stat(upstream)

	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, ErrNotDirectory{}
	}

	repo, err := &Repository{
		Upstream: upstream,
		ValidUpstreams: []string{
			`(.*)`,
		}}, nil

	err = repo.validateUpstream()

	if err != nil {
		return nil, err
	}

	t := configelements.NewTemplDir()
	repo.Destination = fmt.Sprintf("%s/local/%s", t.TemplatesDir, path.Base(repo.Upstream))

	return repo, nil
}

// Fetch implements the CommonRepositoryBehaviour interface for Repository.
func (r Repository) Fetch() error {

	err := fetch(r)

	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Origin() string {
	return r.Upstream
}

func (r Repository) TemplDestination() string {
	return r.Destination
}

// Finally, some private functions
func (r Repository) validateUpstream() error {
	// Iterate through the patterns and attempt to match the URL
	for _, pattern := range r.ValidUpstreams {
		re := regexp.MustCompile(pattern)
		if re.MatchString(r.Upstream) {
			return nil
		}
	}

	return ErrInvalidUpstream{}
}

// No need for this to be on a repository.
func fetch(repo Repository) error {
	return nil
}
