package templcommands

import (
	"github.com/google/subcommands"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestListFiles(t *testing.T) {
	// Set up a test matrix
	testCases := []TestSetup{
		{name: "No argument given", setupFiles: TestFileStructure{directories: []string{}, files: []string{}}, startDirs: []string{"./"}, want: []string{}},

		{name: "Only the top level dir", setupFiles: TestFileStructure{directories: []string{"./"}, files: []string{}}, startDirs: []string{"./"}, want: []string{}},

		{name: "One test file in top dir", setupFiles: TestFileStructure{directories: []string{"./"}, files: []string{"./test1"}}, startDirs: []string{"./"}, want: []string{"test1"}},

		{name: "One test file one dir down", setupFiles: TestFileStructure{directories: []string{"./", "./a_directory"}, files: []string{"./a_directory/test1"}}, startDirs: []string{"./"}, want: []string{"a_directory/test1"}},
	}

	for _, tt := range testCases {

		tempdir := tt.Setup()

		err := os.Chdir(tempdir)
		if err != nil {
			panic(err)
		}

		defer tt.TearDown(tempdir)

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
