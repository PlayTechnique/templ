package templatedirectories

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"templ/configelements"
)

//func Update() []string {
//	startDir := configelements.NewTemplDir().TemplatesDir
//
//}

// FindRepositories searches for all .git directories in TemplatesDir. For each discovered .git directory, it appends
// the containing directory to a list and returns the list.
func FindRepositories() (directories []string, err error) {
	startDir := configelements.NewTemplDir().TemplatesDir

	err = filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && filepath.Base(path) == ".git" {
			// Found a .git directory; add the parent directory to the list.
			repoDir := filepath.Dir(path)
			directories = append(directories, repoDir)
		}

		return nil
	})

	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("%s:%d: %v", file, line, err)
	}

	sort.Strings(directories)

	return directories, nil
}
