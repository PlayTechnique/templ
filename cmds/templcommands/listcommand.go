package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type ListCommand struct {
	TemplateName string
	synopsis     string
	usage        string
}

var listCommand ListCommand

func init() {
	listCommand.TemplateName = "list"
	listCommand.synopsis = "list available templates"
	listCommand.usage = `
list 
Outputs the names of all known template files.
`
}

func (ListCommand) Name() string {
	return listCommand.TemplateName
}

func (ListCommand) Synopsis() string {
	return listCommand.synopsis
}

func (ListCommand) Usage() string {
	return listCommand.usage
}

func (c *ListCommand) SetFlags(_ *flag.FlagSet) {

}

func (*ListCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logrus.Debug(f)
	files, err := listFiles([]string{"."})

	fmt.Println(strings.Join(files, "\n"))

	return err
}

func listFiles(topLevelDirs []string) ([]string, subcommands.ExitStatus) {
	allFileNames := make([]string, 0)

	//For each file named in the args to cat, search the current working directory to see if it exists
	for _, root := range topLevelDirs {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			filename, err := filepath.Rel(root, path)

			if err != nil {
				logrus.Error("Error calculating relative path:", err)
				return err
			}

			// The first file node has a name of '.' for the pwd.
			// It's quicker to check that for skipping than to check the root variable.
			if !info.IsDir() {
				allFileNames = append(allFileNames, filename)
			}
			return nil
		})

		if err != nil {
			logrus.Error(err)
			return nil, subcommands.ExitFailure
		}
	}

	logrus.Debug("Found files: ", allFileNames)

	return allFileNames, subcommands.ExitSuccess
}
