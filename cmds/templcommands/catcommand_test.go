package templcommands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindFilesByNameFindsFilesByName(t *testing.T) {
	testCases := []TestSetup{
		{name: "No argument given", setupFiles: TestFileStructure{directories: []string{}, files: []string{}}, startDirs: []string{"./"}, want: []string{}},

		{name: "Only the top level dir", setupFiles: TestFileStructure{directories: []string{"./"}, files: []string{}}, startDirs: []string{"./"}, want: []string{}},

		{name: "One test file in top dir", setupFiles: TestFileStructure{directories: []string{"./"}, files: []string{"./test1"}}, startDirs: []string{"./"}, want: []string{"test1"}},

		{name: "One test file one dir down", setupFiles: TestFileStructure{directories: []string{"./", "./a_directory"}, files: []string{"./a_directory/test1"}}, startDirs: []string{"./"}, want: []string{"a_directory/test1"}},
	}

	for _, tt := range testCases {
		tempdir := tt.Setup()
		defer tt.TearDown(tempdir)

		filesToFind := makeSet(tt.want)
		// A little prefactoring; if we ever have more than 1 temp dir e.g. when we
		// support multiple template directories, we're already set up for.

		foundFiles, err := findFilesByName(tt.startDirs[0], filesToFind)

		assert.Nil(t, err)
		assert.Equal(t, foundFiles, tt.want)
	}

}
