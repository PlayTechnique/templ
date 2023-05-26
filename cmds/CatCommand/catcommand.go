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
	TemplateName string
	synopsis     string
	usage        string
	mew          bool
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

func (c *CatCommand) SetFlags(_ *flag.FlagSet) {

}

func (*CatCommand) Execute(_ context.Context, _ *flag.FlagSet, subcommandArgs ...interface{}) subcommands.ExitStatus {

	logrus.Debug(subcommandArgs...)

	catArgs := makeSet(subcommandArgs[2:])

	//For each arg in the set, search the current working directory to see if a file exists with it as the name
	// search this directory for a file with the name of arg
	matchingfiles, err := FindFilesByName(".", catArgs)
	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	}

	for _, file := range matchingfiles {
		f, err := os.Open(file)
		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		}
		defer f.Close()

		// Print the contents of the file
		_, err = fmt.Print(f)

		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		}

	}

	return subcommands.ExitSuccess
}

// FindFilesByName searches a directory for file names that match those provided in a set of strings.
// Arguments:
// dir: the directory to search
// names: a set of filenames to search for
// Returns:
// an array of strings, each of which is the path to a file that was found.
// or an error
func FindFilesByName(dir string, names map[string]struct{}) ([]string, error) {
	var foundFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the file's name is in the set of names
		if _, ok := names[info.Name()]; ok {
			foundFiles = append(foundFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return foundFiles, nil
}

// MakeSet creates a set from an array of anything.
// Arguments:
// arr: an array of anything
// Returns:
// a set of strings. Go doesn't have a native set, so use a map with empty structs as the values.
func makeSet(arr []interface{}) map[string]struct{} {
	set := make(map[string]struct{})
	for _, elem := range arr {
		if str, ok := elem.(string); ok {
			set[str] = struct{}{}
		}
	}
	return set
}
