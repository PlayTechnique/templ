package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
)

type RepoCommand struct {
	synopsis           string
	TemplateName       string
	templatedirectory  string
	UpdateRepositories bool
	usage              string
}

var repocommand RepoCommand

func init() {
	repocommand.TemplateName = "repo"
	repocommand.synopsis = "Clone or update git repositories of templates"
	repocommand.usage = `
repo git-url
Clones a git repository into the templates directory. It could be any git repository, but it's best if it contains
templates'
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

func (r RepoCommand) SetFlags(f *flag.FlagSet) {
	// update flag
	f.BoolVar(&repocommand.UpdateRepositories, "update", false, "Update the template repositories")
}

func (r RepoCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logrus.Debug(repocommand.UpdateRepositories)

	if repocommand.UpdateRepositories {
		err := UpdateTheRepositories()
		if err != nil && err != git.NoErrAlreadyUpToDate {
			logrus.Error(err)
			return subcommands.ExitFailure
		} else {
			logrus.Debug("Updated repositories")
		}
		return subcommands.ExitSuccess
	}

	err := repocommand.CloneTheRepositories(f.Args())

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure
	} else {
		logrus.Debug("Cloned repositories from <", f.Args())
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

func UpdateTheRepositories() error {
	repositories, err := getRepositories()

	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, repository := range repositories {
		repo, err := git.PlainOpen(repository)

		if err != nil {
			logrus.Error(err)
			return err
		}

		w, err := repo.Worktree()

		if err != nil {
			logrus.Error(err)
			return err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})

		if err == git.ErrUnstagedChanges {
			logrus.Error("Unstaged changes in ", repository, " not pulling")
			logrus.Error("Repeated calls to update will report success, but this condition will not have changed.")
			logrus.Error("Please manually clean up", repository)
			continue
		}

		if err != nil && err != git.NoErrAlreadyUpToDate {
			err = fmt.Errorf("failed to pull for repo %s : %w", repository, err)
			logrus.Error(err)
			continue
		}

		if err == git.NoErrAlreadyUpToDate {
			logrus.Info(err)
			return err
		}

		logrus.Info("Updated ", repository)
	}

	return nil
}

func getRepositories() ([]string, error) {
	startDir := repocommand.templatedirectory

	// iterate over the directories, looking for the .git directory
	var repositories []string
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".git" {
			// add the parent directory to the list of repositories
			repositories = append(repositories, filepath.Dir(path))
		}
		return nil
	})

	return repositories, err
}

func NewRepoCommand(templateDir string) RepoCommand {
	repocommand.templatedirectory = templateDir
	return repocommand
}
