package templcommands

import (
	"github.com/google/subcommands"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"
)

// Clone a git repository using the `repo` subcommand and then list it using the `list` subcommand

func TestCloneAndList(t *testing.T) {
	topLevel, err := os.MkdirTemp("", "gwyntegrationtest")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(topLevel)

	err = cloneTheRepositories(topLevel, []string{"https://github.com/gwynforthewyn/templ"})
	if err != nil {
		t.Fatal(err)
	}

	files, exitStatus := listFiles([]string{topLevel})

	if exitStatus == subcommands.ExitFailure {
		t.Fatal("List command failed")
	}

	//verify known good strings are in the file list
	search := "templ/main.go"
	info, err := os.Lstat(topLevel + "/templ")

	index := sort.SearchStrings(files, search)

	assert.True(t, info.Mode()&os.ModeSymlink == os.ModeSymlink)
	assert.True(t, index < len(files))
	assert.True(t, files[index] == search)

}
