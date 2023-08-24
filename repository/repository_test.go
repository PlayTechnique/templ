package repository_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"templ/repository"
	"testing"
)

//
//type gitGoMock struct {
//	mock.Mock
//}
//
//func (g gitGoMock) PlainClone {
//
//}

var testLogger *log.Logger
var buf bytes.Buffer

func init() {
	testLogger = log.New(&buf, "repositoryTest: ", log.Lshortfile)
}

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

func TestLocalGitDestination(t *testing.T) {
	tempDir := setupTemplDir("templ-local-github-destination", t)
	defer cleanUpTemplDir(tempDir, t)

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
	defer cleanUpTemplDir(tempDir, t)

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

	testLogger.Printf("repository test dir is %s", tempDir)
	fmt.Print(&buf)

	if err != nil {
		t.Errorf("could not create temp dir")
	}

	err = os.Setenv("TEMPL_DIR", tempDir)

	if err != nil {
		t.Errorf("error setting environment variable TEMPL_DIR")
	}

	return tempDir
}

func cleanUpTemplDir(tempDir string, t *testing.T) {
	err := os.Unsetenv("TEMPL_DIR")

	if err != nil {
		t.Errorf("could not unset TEMPL_DIR env var: %v", err)
	}

	err = os.RemoveAll(tempDir)

	if err != nil {
		t.Errorf("could not remove temp dir: %v", err)
	}
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
