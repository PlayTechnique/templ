package templcommands

import (
	"github.com/google/subcommands"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"testing"
)

type TestFileStructure struct {
	directories []string
	files       []string
}

func TestListFiles(t *testing.T) {
	// Set up a test matrix
	testCases := []struct {
		name       string
		setupFiles TestFileStructure
		startDirs  []string
		want       []string
	}{
		{"No argument given", TestFileStructure{directories: []string{}}, []string{"./"}, []string{}},
		{"Only the top level dir", TestFileStructure{directories: []string{"./"}}, []string{"./"}, []string{}},
		{"One test file in top dir", TestFileStructure{directories: []string{"./"}, files: []string{"./test1"}}, []string{"./"}, []string{"test1"}},
		{"One test file one dir down", TestFileStructure{directories: []string{"./", "./a_directory"}, files: []string{"./a_directory/test1"}}, []string{"./"}, []string{"a_directory/test1"}},
	}

	for _, tt := range testCases {
		tempdir := setup(tt.setupFiles)

		err := os.Chdir(tempdir)
		if err != nil {
			panic(err)
		}

		defer os.RemoveAll(tempdir)

		t.Run(tt.name, func(t *testing.T) {
			ans, err := listFiles(tt.startDirs)

			assert.True(t, err == subcommands.ExitSuccess, "Expected success, got %s", err)

			sort.Strings(ans)
			sort.Strings(tt.want)
			assert.True(t, reflect.DeepEqual(ans, tt.want), "Expected %s, got %s", tt.want, ans)

		},
		)
	}
}

func setup(structure TestFileStructure) (tempdir string) {
	tempdir, err := os.MkdirTemp("", "templ_test")

	if err != nil {
		panic(err)
	}

	err = os.Chdir(tempdir)
	if err != nil {
		panic(err)
	}

	for _, dir := range structure.directories {
		os.MkdirAll(dir, 0755)
	}

	for _, files := range structure.files {
		os.Create(files)
	}

	return tempdir
}
