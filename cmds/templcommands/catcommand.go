package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type CatCommand struct {
	TemplateName      string
	templatedirectory string
	synopsis          string
	usage             string
}

var catcommand CatCommand

func init() {
	catcommand.TemplateName = "cat"
	catcommand.synopsis = "cat a named template templates"
	catcommand.usage = `
cat <name of a template file> 
Output the contents of named, known template files.

You can see the names of the template files with the list subcommand.
`
}

func (CatCommand) Name() string {
	return catcommand.TemplateName
}

func (CatCommand) Synopsis() string {
	return catcommand.synopsis
}

func (CatCommand) Usage() string {
	return catcommand.usage
}

func (c CatCommand) SetFlags(_ *flag.FlagSet) {

}

func (c CatCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logrus.Debug(f)

	//For each file named in the args to cat, search the current working directory to see if it exists.
	matchingFiles, err := findFilesByName(c.templatedirectory, f.Args())
	logrus.Debug("Found files: ", matchingFiles)

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	}

	for _, file := range matchingFiles {
		contents, err := os.ReadFile(file)

		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		}

		// Print the contents of the file
		_, err = fmt.Print(string(contents))

		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		}

	}

	return subcommands.ExitSuccess
}

func NewCatCommand(templatedirectory string) CatCommand {
	catcommand.templatedirectory = templatedirectory
	return catcommand
}

func getAbsPaths(paths []string, root string) ([]string, error) {
	var absPaths = make([]string, 0)

	for _, path := range paths {
		absPath, err := filepath.Abs(root + "/" + path)

		if err != nil {
			return absPaths, err
		}

		absPaths = append(absPaths, absPath)
	}

	return absPaths, nil
}
