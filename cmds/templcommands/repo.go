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
	TemplateName          string
	templatedirectory     string
	synopsis              string
	usage                 string
	repositoriesdirectory string
}

var repocommand RepoCommand

func init() {
	repocommand.TemplateName = "repo"
	repocommand.synopsis = "Clone or update git repositories of templates"
	repocommand.usage = `
repo git-url
Clones a github repository

repo --update
Iterates over repositories and updates them
`
}

func (RepoCommand) Name() string {
	return repocommand.TemplateName
}

func (RepoCommand) Synopsis() string {
	return repocommand.synopsis
}

func (RepoCommand) Usage() string {
	return repocommand.usage
}

func (c RepoCommand) SetFlags(_ *flag.FlagSet) {

}

func (r RepoCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logrus.Debug(f)

	if f.NArg() == 0 {
		fmt.Print(repocommand.Usage())
		return subcommands.ExitFailure
	}

	err := repocommand.CloneTheRepositories(f.Args())

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	} else {
		logrus.Debug("Cloned repositories from <", f.Args(), "> into the ", repocommand.repositoriesdirectory, " directory.")
	}

	return subcommands.ExitSuccess
}

func (r RepoCommand) CloneTheRepositories(gitUrls []string) error {

	for _, gitUrl := range gitUrls {
		destinationDirectoryName, err := extractDirNameFromUrl(gitUrl)

		if err != nil {
			logrus.Error("Could not extract directory name from ", gitUrl, " ", err)
			return err
		}

		_, err = git.PlainClone(r.templatedirectory+"/"+destinationDirectoryName, false, &git.CloneOptions{
			URL:      gitUrl,
			Progress: os.Stdout,
		})

		if err != nil {
			if err == git.ErrRepositoryAlreadyExists {
				logrus.Info("Repository ", gitUrl, " already exists, not cloning: ", err)
				continue
			}
			logrus.Error(err)
			return err
		}
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

func NewRepoCommand(templateDir string) RepoCommand {
	repocommand.templatedirectory = templateDir
	repocommand.repositoriesdirectory = templateDir + "/.repositories"
	return repocommand
}
