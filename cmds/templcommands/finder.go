package templcommands

import (
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// findFilesByName searches a directory for file names that match those provided in a set of strings.
// Arguments:
// root: the directory to search
// names: a set of filenames to search for
// Returns:
// an array of strings, each of which is the path to a file that was found.
// or an error
func findFilesByName(root string, names map[string]struct{}) ([]string, error) {
	foundFiles := []string{}

	logrus.Debug("Outside filepath.Walk function names: ", names)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logrus.Debug("Inside filepath.Walk function names: ", names)

		// If the file's name is in the set of names
		if _, ok := names[path]; ok {
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

			if !info.IsDir() {
				logrus.Debug("Appending: ", filename)
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
