package templcommands

import (
	"context"
	"flag"
	"github.com/go-git/go-git/v5"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
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

	if len(f.Args()) > 0 {
		err := cloneTheRepositories(repoCommand.repodirectory, f.Args())
		if err != nil {
			logrus.Error(err)
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}

func cloneTheRepositories(repodirectory string, gitUrls []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error("Could not get pwd from", pwd, " ", err)
		return err
	}

	defer os.Chdir(pwd)
	os.Chdir(repodirectory)

	for _, gitUrl := range gitUrls {

		_, err := git.PlainClone(repodirectory, false, &git.CloneOptions{
			URL:      gitUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}

		// symlink the new directory to the templates directory
	}

	return nil
}
