package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	templcommands "playtechnique.io/templ/cmds/CatCommand"
)

func init() {
	lvl, ok := os.LookupEnv("TEMPL_LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "error"
	}

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set file names and line number to appear in log messages
	logrus.SetReportCaller(true)

	// set global log level
	logrus.SetLevel(ll)

}

// Main identifies the templates directory, switches working directories into it and invokes the subcommand.
func main() {
	ctx := context.Background()

	logrus.Info("Starting templ")

	configDir := getConfigDirectory()
	ChangeDir(configDir)

	oldHelp := subcommands.DefaultCommander.Explain

	help := func(w io.Writer) {
		fmt.Fprintf(w, "Env Vars: LOG_LEVEL TEMPL_DIR\n")
		oldHelp(w)

	}
	subcommands.DefaultCommander.Explain = help

	subcommands.Register(subcommands.HelpCommand(), "help")
	subcommands.Register(&templcommands.CatCommand{}, "templates")

	// Mystical. This seems to parse the subcommand flags.
	flag.Parse()
	os.Exit(int(subcommands.Execute(ctx, os.Args)))
}

// Interrogates the environment variable TEMPL_DIR for a directory of templates.
// If that variable does not exist, it attempts to switch into the default template directory.
func getConfigDirectory() string {

	configDir, foundEnvVar := os.LookupEnv("TEMPL_DIR")

	if foundEnvVar {
		return configDir
	}

	home, foundHome := os.LookupEnv("HOME")

	if !foundHome {
		fmt.Print("HOME env var not found. templ cannot autodiscover the default configdir, and TEMPL_DIR is not set.")
		panic("Please review the help for preconditions to meet running templ/")
	}

	configDir = home + "/.config/templ/"

	if _, err := os.Stat(configDir); err == nil || os.IsNotExist(err) {
		fmt.Print("TEMPL_DIR not set. Did not find the default directory. Creating...")
		err := os.Mkdir(configDir, 0700)

		if err != nil {
			panic(err)
		}
	}

	return configDir

}

// Validates that a directory exists and then changes the pwd into it.
// Parameters:
// 1. templatedir - a string representing the path to the directory containing templates.
// Returns:
// string - Absolute path to the directory that has been switched into.
// Side Effects
// Will panic if the directory does not exist.
func ChangeDir(templatedir string) string {
	stat, err := os.Stat(templatedir)

	if !stat.IsDir() || err != nil {
		panic(templatedir + " does not exist")
	} else {
		templatedir, err = filepath.Abs(templatedir)
		if err != nil {
			panic(fmt.Errorf("something wrong with %s: %v", templatedir, err))
		}
		os.Chdir(templatedir)
		return templatedir
	}

}
