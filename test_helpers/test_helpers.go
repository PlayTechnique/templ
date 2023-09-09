package test_helpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"runtime"
	"templ/repository"
	"testing"
)

func CleanUpTemplDir(tempDir string, t *testing.T) {
	err := os.Unsetenv("TEMPL_DIR")

	if err != nil {
		t.Errorf("could not unset TEMPL_DIR env var: %v", err)
	}

	err = os.RemoveAll(tempDir)

	if err != nil {
		t.Errorf("could not remove temp dir: %v", err)
	}
}

// CreateTemplDirWithMultipleRepositories will clone the testing-files git submodule's git repository
// into several directories named "repositories" in the second.
func CreateTemplDirWithMultipleRepositories(cloneDestinations []string) (tempDir string, repositories []string, err error) {
	testingFiles, err := filepath.Abs("../testing-files")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", repositories, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	tempDir, err = os.MkdirTemp("", "testing-git-clones")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return tempDir, repositories, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	repoToCloneFrom, err := repository.NewGitRepository(testingFiles)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", repositories, fmt.Errorf("%s:%d: %w", file, line, err)
	}

	for _, dest := range cloneDestinations {

		repoDest := tempDir + "/" + dest

		dest, err := git.PlainClone(repoDest, false, &git.CloneOptions{
			URL:               repoToCloneFrom.Origin(),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return tempDir, repositories, fmt.Errorf("%s:%d: %v:", file, line, err)
		}

		fmt.Printf("%v", dest)
		repositories = append(repositories, repoDest)
	}

	return
}
