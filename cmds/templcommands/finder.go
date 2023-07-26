package templcommands

import (
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// findFilesByName searches a directory for file names that match those provided in a set of strings.
// Arguments:
// root: the directory to search
// names: a set of filenames to search for
// Returns:
// an array of strings, each of which is the path to a file that was found.
// or an error
func findFilesByName(root string, names []string) ([]string, error) {
	foundFiles := []string{}

	logrus.Debug("Outside filepath.Walk function names: ", names)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logrus.Debug("Inside filepath.Walk function names: ", names)

		// If the file's name is in the set of names
		for _, name := range names {
			logrus.Debug("name is <", name, "> path is <", path, ">")
			if strings.Contains(path, name) {
				foundFiles = append(foundFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return foundFiles, nil
}

func listFiles(topLevelDirs []string) ([]string, subcommands.ExitStatus) {
	allFileNames := make([]string, 0)

	//For each file named in the args to cat, search the current working directory to see if it exists
	for _, root := range topLevelDirs {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			filename, err := filepath.Rel(root, path)

			// If we're at the root, don't add it to the list of files, but also don't let it get to
			// the filepath.Skipdir return statement, otherwise we don't walk anything at all.
			if path == root {
				return nil
			}

			// skip hidden files and directories
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

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
