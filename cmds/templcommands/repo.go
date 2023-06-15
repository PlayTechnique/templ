package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
)

type RepoCommand struct {
	TemplateName  string
	synopsis      string
	usage         string
	repodirectory string
}

var repoCommand RepoCommand

func init() {
	repoCommand.TemplateName = "repo"
	repoCommand.synopsis = "Clone or update git repositories of templates"
	repoCommand.usage = `
repo git-url
Clones a github repository

repo --update
Iterates over repositories and updates them

`
	// We'll make a directory called .repositories in the current working directory
	// to store any git repositories that are cloned down.
	// It feels simpler than cloning the directories directly into the templates directory and then trying to
	// differentiate between user-provided repositories and the ones that are cloned down.
	repoCommand.repodirectory = ".repositories"
}

func (RepoCommand) Name() string {
	return repoCommand.TemplateName
}

func (RepoCommand) Synopsis() string {
	return repoCommand.synopsis
}

func (RepoCommand) Usage() string {
	return repoCommand.usage
}

func (c *RepoCommand) SetFlags(_ *flag.FlagSet) {

}

func (*RepoCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logrus.Debug(f)

	if f.NArg() == 0 {
		fmt.Print(repoCommand.Usage())
		return subcommands.ExitFailure
	}

	//if r.updates {
	//	doTheRepoUpdates()
	//}

	// if repodirectory does not exist, create it
	if _, err := os.Stat(repoCommand.repodirectory); os.IsNotExist(err) {
		err := os.MkdirAll(repoCommand.repodirectory, os.ModePerm)

		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		} else {
			logrus.Debug("Created directory: ", repoCommand.repodirectory)
		}
	}

	err := cloneTheRepositories(repoCommand.repodirectory, f.Args())

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	} else {
		logrus.Debug("Cloned repositories from <", f.Args(), "> into the ", repoCommand.repodirectory, " directory.")
	}

	return subcommands.ExitSuccess
}

func cloneTheRepositories(repositoriesTopLevel string, gitUrls []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error("Could not get pwd from", pwd, " ", err)
		return err
	}

	for _, gitUrl := range gitUrls {

		destinationDirectoryName, err := extractDirNameFromUrl(gitUrl)

		if err != nil {
			logrus.Error("Could not extract directory name from ", gitUrl, " ", err)
			return err
		}

		_, err = git.PlainClone(repositoriesTopLevel+"/"+destinationDirectoryName, false, &git.CloneOptions{
			URL:      gitUrl,
			Progress: os.Stdout,
		})

		if err != nil {
			if err == git.ErrRepositoryAlreadyExists {
				logrus.Info("Repository ", gitUrl, " already exists, not cloning: ", err)
			}
			logrus.Error(err)
			return err
		}

		// symlink the new directory to the templates directory
		err = os.Symlink(repositoriesTopLevel+"/"+destinationDirectoryName, destinationDirectoryName)
	}

	return nil
}

func extractDirNameFromUrl(urlStr string) (string, error) {
	// Regular expression pattern for GitHub repository extraction
	patterns := []string{
		`github\.com[:/]([^/]+)/([^/.]+)(\.git)?$`,        // git@github.com:<owner>/<repo> or https://github.com/<owner>/<repo>
		`https?://github\.com/([^/]+)/([^/.]+)(\.git)?$`,  // https://github.com/<owner>/<repo>
		`git://github\.com/([^/]+)/([^/.]+)(\.git)?$`,     // git://github.com/<owner>/<repo>
		`ssh://git@github\.com/([^/]+)/([^/.]+)(\.git)?$`, // ssh://git@github.com/<owner>/<repo>
		`git@github\.com:([^/]+)/([^/.]+)(\.git)?$`,       // git@github.com:<owner>/<repo>.git
	}

	// Iterate through the patterns and attempt to match the URL
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(urlStr)
		if matches != nil {
			// Return the repository name
			repository_name := matches[2]
			repo := repository_name
			return repo, nil
		}
	}

	return "", fmt.Errorf("unable to extract repository from URL")
}
