package templcommands

import (
	"github.com/google/subcommands"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestRenderFile(t *testing.T) {
	// Set up a test matrix
	testCases := []TestSetup{
		{name: "No argument given", setupFiles: TestFileStructure{directories: []string{}, files: []string{}}, startDirs: []string{"./"}, want: []string{}},
	}

	for _, tt := range testCases {

		tempdir := Setup(tt.setupFiles)

		err := os.Chdir(tempdir)
		if err != nil {
			panic(err)
		}

		defer TearDown(tempdir)

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

func TestFindFilesInArgs(t *testing.T) {
	setupStructure := TestFileStructure{directories: []string{"level_one", "level_one/level_two"}, files: []string{"ringding", "level_one/smudge"}}

	tempdir := Setup(setupStructure)
	err := os.Chdir(tempdir)
	if err != nil {
		panic(err)
	}

	defer TearDown(tempdir)

	fin
}
