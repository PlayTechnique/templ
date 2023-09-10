package templatedirectories_test

import (
	"reflect"
	"sort"
	"templ/templatedirectories"
	"templ/test_helpers"
	"testing"
)

func TestListWithNoDirectoriesCloned(t *testing.T) {
	_, err := test_helpers.CreateFileSystem([]string{})
	if err != nil {
		t.Errorf("%v", err)
	}

	directories, err := templatedirectories.List()

	if err != nil {
		t.Errorf("Could not list directories: %v", err)
	}

	expectedDirectories := []string{}

	if !reflect.DeepEqual(directories, expectedDirectories) {
		t.Errorf("Expected <%v>, got <%v>", expectedDirectories, directories)
	}
}

func TestListDoesNotFindHiddenDirectories(t *testing.T) {
	create := []string{".hidden/"}
	templDir, err := test_helpers.CreateFileSystem(create)
	defer test_helpers.CleanUpTemplDir(templDir, t)

	if err != nil {
		t.Errorf("%v", err)
	}

	directories, err := templatedirectories.List()

	if err != nil {
		t.Errorf("%v", err)
	}

	expectedDirectories := []string{}
	if !reflect.DeepEqual(directories, expectedDirectories) {
		t.Errorf("Expected <%v>, got <%v>", expectedDirectories, directories)
	}
}

func TestListWithOneFile(t *testing.T) {
	createFiles := []string{"file"}
	templDir, err := test_helpers.CreateFileSystem(createFiles)

	if err != nil {
		t.Errorf("%v", err)
	}

	defer test_helpers.CleanUpTemplDir(templDir, t)

	files, err := templatedirectories.List()

	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(files, createFiles) {
		t.Errorf("Expected <%v>, got <%v>", createFiles, files)
	}
}

func TestListWithTemplateFiles(t *testing.T) {
	createFiles := []string{"file", "diggle/flibble/foo.yaml"}
	sort.Strings(createFiles)

	templDir, err := test_helpers.CreateFileSystem(createFiles)

	if err != nil {
		t.Errorf("%v", err)
	}

	defer test_helpers.CleanUpTemplDir(templDir, t)

	files, err := templatedirectories.List()

	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(files, createFiles) {
		t.Errorf("Expected <%v>, got <%v>", createFiles, files)
	}
}
