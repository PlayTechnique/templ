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

func (*CatCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logrus.Debug(f)
	catArgs := makeSet(f.Args())
	logrus.Debug("catArgs: ", catArgs)

	//For each file named in the args to cat, search the current working directory to see if it exists
	matchingfiles, err := FindFilesByName(".", catArgs)
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

// FindFilesByName searches a directory for file names that match those provided in a set of strings.
// Arguments:
// dir: the directory to search
// names: a set of filenames to search for
// Returns:
// an array of strings, each of which is the path to a file that was found.
// or an error
func FindFilesByName(dir string, names map[string]struct{}) ([]string, error) {
	var foundFiles []string
	logrus.Debug("Outside filepath.Walk function names: ", names)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logrus.Debug("Inside filepath.Walk function names: ", names)

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
func makeSet(arr []string) map[string]struct{} {

	set := make(map[string]struct{})

	for _, filename := range arr {
		logrus.Debug("filename: ", filename)
		set[filename] = struct{}{}
	}

	return set
}
