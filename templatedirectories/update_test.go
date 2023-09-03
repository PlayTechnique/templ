package templatedirectories_test

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"templ/repository"
	"templ/templatedirectories"
	"templ/test_helpers"
	"testing"
)

func TestFindingGitRepositories(t *testing.T) {
	templDir, expectedRepositories, err := createTemplDirWithMultipleRepositories([]string{"templ-update-test", "roflcopter/in-here"})
	sort.Strings(expectedRepositories)

	if err != nil {
		t.Errorf("%v", err)
	}

	defer test_helpers.CleanUpTemplDir(templDir, t)
	err = os.Setenv("TEMPL_DIR", templDir)

	if err != nil {
		t.Errorf("%v", err)
	}

	repositories, err := templatedirectories.FindRepositories()
	if err != nil {
		t.Errorf("FindRepositories returned %v", err)
	}

	if !reflect.DeepEqual(repositories, expectedRepositories) {
		t.Errorf("Discovered repositories <%s> do not match expected repositories <%s>", repositories, expectedRepositories)
	}
}

func TestFindingNoGitRepositories(t *testing.T) {
	templDir, expectedRepositories, err := createTemplDirWithMultipleRepositories([]string{})
	sort.Strings(expectedRepositories)

	if err != nil {
		t.Errorf("%v", err)
	}

	defer test_helpers.CleanUpTemplDir(templDir, t)
	err = os.Setenv("TEMPL_DIR", templDir)

	if err != nil {
		t.Errorf("%v", err)
	}

	repositories, err := templatedirectories.FindRepositories()

	if repositories != nil {
		t.Errorf("Should have found no template directories, found %v", repositories)
	}
}

// createTemplDirWithMultipleRepositories will clone the testing-files git submodule's git repository
// into several directories named "repositories" in the second.
func createTemplDirWithMultipleRepositories(cloneDestinations []string) (tempDir string, repositories []string, err error) {
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
