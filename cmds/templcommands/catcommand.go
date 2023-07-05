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
	// If a file name is supplied twice, only have it once in the data structure. This keeps the cat output
	// clean, rather than having duplicated files.
	// The key in the set is the relative path to the file from the os.args array. The value in the set is
	// empty.
	absPaths, err := getAbsPaths(f.Args(), c.templatedirectory)
	catArgs := makeSet(absPaths)
	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	}

	logrus.Debug("catArgs: ", catArgs)

	//For each file named in the args to cat, search the current working directory to see if it exists
	// TODO: change the signature of findFilesByName to take a slice of strings
	matchingfiles, err := findFilesByName(c.templatedirectory, catArgs)
	logrus.Debug("Found files: ", matchingfiles)

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	}

	for _, file := range matchingfiles {
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
