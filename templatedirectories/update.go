package templatedirectories

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"templ/configelements"
)

func Update() error {
	repositories, err := FindRepositories()

	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("%s:%d: %v", file, line, err)
	}

	for _, repository := range repositories {
		repo, err := git.PlainOpen(repository)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			newerr := fmt.Errorf("%s:%d: %v", file, line, err)
			logrus.Info(newerr)
			return newerr
		}

		w, err := repo.Worktree()

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			newerr := fmt.Errorf("%s:%d: %v", file, line, err)
			logrus.Info(newerr)
			return newerr
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})

		if errors.Is(err, git.ErrUnstagedChanges) {
			logrus.Error("Unstaged changes in ", repository, " not pulling")
			logrus.Error("Repeated calls to update will report success, but this condition will not have changed.")
			logrus.Error("Please manually clean up", repository)
			continue
		}

		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				newerr := fmt.Errorf("%s:%d: failed to pull for repo %s: %v", file, line, repository, err)
				logrus.Error(newerr)
				continue
			}
		}

		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logrus.Info(err)
		}

		logrus.Info("Updated ", repository)
	}

	return nil
}

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
