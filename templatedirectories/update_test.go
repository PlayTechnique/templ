package templatedirectories_test

import (
	"os"
	"reflect"
	"sort"
	"templ/templatedirectories"
	"templ/test_helpers"
	"testing"
)

func TestFindingGitRepositories(t *testing.T) {
	templDir, expectedRepositories, err := test_helpers.CreateTemplDirWithMultipleRepositories([]string{"templ-update-test", "roflcopter/in-here"})

	if err != nil {
		t.Errorf("%v", err)
	}

	sort.Strings(expectedRepositories)

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
	templDir, expectedRepositories, err := test_helpers.CreateTemplDirWithMultipleRepositories([]string{})
	sort.Strings(expectedRepositories)

	if err != nil {
		t.Fatal(err)
	}

	defer test_helpers.CleanUpTemplDir(templDir, t)
	err = os.Setenv("TEMPL_DIR", templDir)

	if err != nil {
		t.Fatal(err)
	}

	repositories, err := templatedirectories.FindRepositories()

	if err != nil {
		t.Fatal(err)
	}

	if repositories != nil {
		t.Errorf("Should have found no template directories, found %v", repositories)
	}
}
