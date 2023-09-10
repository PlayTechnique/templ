package repository_test

import (
	"errors"
	"fmt"
	"os"
	"templ/repository"
	"templ/test_helpers"
	"testing"
)

func TestRepositoryConstructorWithEmptyUrl(t *testing.T) {
	_, err := repository.NewGitRepository("")

	if !errors.Is(err, repository.ErrInvalidUpstream{}) {
		t.Errorf("expected error of type ErrInvalidUpstream, got: %v", err)
	}
}

func TestGithubConstructorWithEmptyUrl(t *testing.T) {
	_, err := repository.NewGitHubRepository("")

	if !errors.Is(err, repository.ErrInvalidUpstream{}) {
		t.Errorf("expected error of type ErrInvalidUpstream, got: %v", err)
	}
}

func TestLocalGitConstructorWithEmptyUrl(t *testing.T) {
	_, err := repository.NewLocalGitRepository("")

	if !errors.Is(err, repository.ErrInvalidUpstream{}) {
		t.Errorf("expected error of type ErrInvalidUpstream, got: %v", err)
	}
}

func TestGithubConstructorWithDirectory(t *testing.T) {
	upstream := "../definitely-not-a-directory"
	_, err := repository.NewGitHubRepository(upstream)

	if !errors.Is(err, repository.ErrInvalidUpstream{}) {
		t.Errorf("A github repository should not validate against a non-existent directory. Received %v", err)
	}

}

func TestLocalGitDestination(t *testing.T) {
	tempDir := setupTemplDir("templ-local-github-destination", t)
	defer test_helpers.CleanUpTemplDir(tempDir, t)

	// The current working directory is the directory repository_test.go is in.
	localTestRepo := "../testing-files/local-git-repo"
	destDir := fmt.Sprintf("%s/local/local-git-repo", tempDir)

	gitRepo, err := repository.NewLocalGitRepository(localTestRepo)
	if err != nil {
		t.Errorf("could not find %s: %v", localTestRepo, err)
	}

	if gitRepo.TemplDestination() != destDir {
		t.Errorf("local repo templ destination is %s, should be %s", gitRepo.TemplDestination(), destDir)
	}

}

func TestGithubDestination(t *testing.T) {
	tempDir := setupTemplDir("templ-local-github-destination", t)
	defer test_helpers.CleanUpTemplDir(tempDir, t)

	// The current working directory is the directory repository_test.go is in.
	destDir := fmt.Sprintf("%s/github/PlayTechnique/templ_templates", tempDir)

	gitRepo, err := repository.NewGitHubRepository("https://github.com/PlayTechnique/templ_templates.git")

	if err != nil {
		t.Errorf("error constructing GithubRepositoryfrom %s: %v", gitRepo, err)
	}

	if gitRepo.TemplDestination() != destDir {
		t.Errorf("github repo templ destination is %s, should be %s", gitRepo.TemplDestination(), destDir)
	}

}

// Creates a temporary directory and sets the TEMPL_DIR environment variable.
func setupTemplDir(pattern string, t *testing.T) string {
	tempDir, err := os.MkdirTemp("", pattern)

	if err != nil {
		t.Errorf("could not create temp dir")
	}

	err = os.Setenv("TEMPL_DIR", tempDir)

	if err != nil {
		t.Errorf("error setting environment variable TEMPL_DIR")
	}

	return tempDir
}

func TestGithubFetchWithInvalidUrl(t *testing.T) {
	_, err := repository.NewGitRepository("../definitely-not-real/")

	if !errors.Is(err, repository.ErrInvalidUpstream{}) {
		t.Errorf("Unexpected error testing git fetch. Expected ErrInvalidUpstream, got %v", err)
	}

}

//func TestGithubFetchWithValidUrl(t *testing.T) {
//
//}
