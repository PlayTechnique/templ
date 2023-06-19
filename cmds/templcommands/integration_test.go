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

	var repocommand RepoCommand = NewRepoCommand(topLevel)

	err = repocommand.CloneTheRepositories([]string{"https://github.com/gwynforthewyn/templ"})
	if err != nil {
		t.Fatal(err)
	}

	filesInTemplateDirectory, exitStatus := listFiles([]string{repocommand.templatedirectory})

	if exitStatus == subcommands.ExitFailure {
		t.Fatal("List command failed")
	}

	repositoryDir, err := os.Stat(topLevel + "/templ")
	index := sort.SearchStrings(filesInTemplateDirectory, "templ/main.go")

	//assert that the templ directory exists
	assert.True(t, repositoryDir.Mode()&os.ModeDir == os.ModeDir)
	assert.True(t, filesInTemplateDirectory[index] == "templ/main.go")

}
