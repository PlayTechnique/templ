package test_helpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func CreateFileSystem(paths []string) (tempDir string, err error) {
	tempDir, err = os.MkdirTemp("", "testing-templ")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	err = os.Setenv("TEMPL_DIR", tempDir)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	for _, p := range paths {
		fullPath := tempDir + "/" + p

		// If the path ends with a "/", it's a directory.
		if strings.HasSuffix(fullPath, "/") {
			err := os.MkdirAll(fullPath, 0755) // 0755 is a common permission setting for directories
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
			}
		} else { // It's a path to a file

			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
			}

			// Create an empty file
			file, err := os.Create(fullPath)
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
			}
			err = file.Close()
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				return tempDir, fmt.Errorf("%s:%d: %v", file, line, err)
			}
		}
	}
	return tempDir, err
}
