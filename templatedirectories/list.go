package templatedirectories

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"templ/configelements"
)

// List lists the template files in the templates directory.
// It does not descend into hidden directories; it does not return hidden files;
func List() ([]string, error) {
	allFileNames := []string{}
	templDir := configelements.NewTemplDir().TemplatesDir

	err := filepath.Walk(templDir, func(path string, info os.FileInfo, err error) error {
		filename, err := filepath.Rel(templDir, path)

		if err != nil {
			return err
		}

		// If we're at the root, don't add it to the list of files, but also don't let it get to
		// the filepath.Skipdir return statement, otherwise we don't walk anything at all.
		if path == templDir {
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
			allFileNames = append(allFileNames, filename)
		}

		return nil
	})

	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return allFileNames, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	sort.Strings(allFileNames)

	return allFileNames, nil
}
