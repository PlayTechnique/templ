package configelements

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// custom errors
type ErrTemplDirPrecondition struct {
}

func (e ErrTemplDirPrecondition) Error() string {
	return "Please review the help for preconditions to meet for running templ"
}

type TemplDir struct {
	TemplatesDir string
	DefaultDir   string
	EnvVar       string
}

func NewTemplDir() TemplDir {
	t := TemplDir{EnvVar: "TEMPL_DIR", DefaultDir: templDirDefault()}

	t.TemplatesDir = t.getTemplatesDirectory()

	return t
}

func templDirDefault() string {
	home, foundHome := os.LookupEnv("HOME")

	if !foundHome {
		logrus.Error("HOME env var not found. templ cannot autodiscover the default templates dir, and TEMPL_DIR is not set.")
		panic(ErrTemplDirPrecondition{})
	}

	return home + "/.config/templ"
}

// GetTemplatesDirectory interrogates the environment variable TEMPL_DIR for a directory of templates.
// If that variable does not exist, we use the default template directory.
// This function will also create the templates directory if needed.
// Note that this function is so foundational to this program that it panics on error rather than
// returning an error value.
func (t TemplDir) getTemplatesDirectory() string {

	var templatesDir string
	var foundEnvVar bool
	var err error

	templatesDir, foundEnvVar = os.LookupEnv(t.EnvVar)

	// Only auto-create TEMPL_DIR if we're using the default directory. The reasoning is that if
	// the user explicitly provides the directory to use then they're expected to have already
	// created it. If not, then we are make sure the default setup is correct.
	if foundEnvVar {
		logrus.Debug("Found ", t.EnvVar, "=", templatesDir)
	} else {
		logrus.Debug("Did not find", t.EnvVar, "Switching to default templates directory.")
		templatesDir = t.DefaultDir
	}

	_, err = os.Stat(templatesDir)

	if os.IsNotExist(err) {
		logrus.Info("Did not find <" + templatesDir + "> directory. Creating...")
		err := os.MkdirAll(templatesDir, 0700)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	templatesDir, err = filepath.Abs(templatesDir)

	if err != nil {
		err = fmt.Errorf("could not take absolute path of %s: %v", templatesDir, err)
		panic(err)
	}

	return templatesDir
}
