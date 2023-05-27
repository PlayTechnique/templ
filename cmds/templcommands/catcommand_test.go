package templcommands

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type TestSetup struct {
	Files map[string]struct{}
}

// Creates a temp directory and populates it with files.
func (t TestSetup) SetupTests() (tempDir string) {

	// Create a directory with some files in it
	tempDir, err := os.MkdirTemp("", "templ_test")

	if err != nil {
		panic(err)
	}

	err = os.Chdir(tempDir)

	if err != nil {
		panic(err)
	}

	// Create a directory with some files in it
	for filename := range t.Files {
		destFile := tempDir + filename
		_, err := os.Create(destFile)

		if err != nil {
			panic(err)
		}
	}

	return tempDir
}

func TearDown(directories []string) {

	os.Chdir("/tmp")
	for _, directory := range directories {
		err := os.RemoveAll(directory)
		if err != nil {
			panic(err)
		}
	}

}

func TestFindFilesByNameFindsFilesByName(t *testing.T) {
	testFiles := map[string]struct{}{
		"test1": {},
		"test2": {},
		"test3": {},
	}

	var testInfo = TestSetup{Files: testFiles}

	tempDir := testInfo.SetupTests()

	// A little prefactoring; if we ever have more than 1 temp dir e.g. when we
	// support multiple template directories, we're already set up for.
	collectionOfTestDir := []string{tempDir}
	defer TearDown(collectionOfTestDir)

	foundFiles, err := FindFilesByName(".", testFiles)

	assert.Nil(t, err)
	assert.Equal(t, len(testFiles), len(foundFiles))

}
