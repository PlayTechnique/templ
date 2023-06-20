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
	"github.com/gwynforthewyn/templ/cmds/templcommands"
)

func init() {
	debugger, ok := os.LookupEnv("TEMPL_DEBUG_BREAK")

	if ok && debugger != "" {
		fmt.Println("Attach debugger to PID: ", os.Getpid())
		fmt.Println("Press enter to continue")
		fmt.Scanln()

	}

	lvl, ok := os.LookupEnv("TEMPL_LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "warn"
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
	tflags := flag.NewFlagSet("templFlags", flag.ExitOnError)

	templCommander := subcommands.NewCommander(tflags, os.Args[0])

	logrus.Info("Starting ", os.Args[0])

	templatesDir, err := getTemplatesDirectory()
	if err != nil {
		panic(err)
	}

	//Super cheap way to get documentation into the usage message.
	oldHelp := templCommander.Explain
	help := func(w io.Writer) {
		fmt.Fprintf(w, "Env Vars: TEMPL_LOG_LEVEL TEMPL_DIR\n")
		oldHelp(w)

	}

	templCommander.Explain = help

	templCommander.Register(subcommands.HelpCommand(), "help")
	templCommander.Register(templcommands.NewCatCommand(templatesDir), "templates")
	templCommander.Register(templcommands.NewListCommand(templatesDir), "templates")
	templCommander.Register(templcommands.NewRepoCommand(templatesDir), "templates")
	templCommander.Register(templcommands.NewRenderCommand(templatesDir), "templates")

	tflags.Parse(os.Args[1:])

	exitval := int(templCommander.Execute(ctx, os.Args))
	os.Exit(exitval)
}

// Interrogates the environment variable TEMPL_DIR for a directory of templates.
// If that variable does not exist, it attempts to switch into the default template directory.
func getTemplatesDirectory() (string, error) {

	templatesDir, foundEnvVar := os.LookupEnv("TEMPL_DIR")

	// Only auto-create TEMPL_DIR if we're using the default directory. The reasoning is that if
	// the user explicitly provides the directory to use then they're expected to have already
	// created it. If not, then we are make sure the default setup is correct.
	if foundEnvVar {
		logrus.Debug("Found TEMPL_DIR=", templatesDir)
	} else {
		logrus.Debug("Did not find TEMPL_DIR. Switching to default templates directory.")
		home, foundHome := os.LookupEnv("HOME")

		if !foundHome {
			logrus.Error("HOME env var not found. templ cannot autodiscover the default templates dir, and TEMPL_DIR is not set.")
			panic("Please review the help for preconditions to meet running templ/")
		}

		templatesDir = home + "/.config/templ"

		if _, err := os.Stat(templatesDir); err == nil || os.IsNotExist(err) {
			logrus.Info("Did not find <" + templatesDir + "> directory. Creating...")
			err := os.MkdirAll(templatesDir, 0700)

			if err != nil {
				panic(err)
			}
		}
	}

	templatesDir, err := filepath.Abs(templatesDir)
	if err != nil {
		err = fmt.Errorf("something wrong with %s: %v", templatesDir, err)
	}

	return templatesDir, err

}
