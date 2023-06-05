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
	"playtechnique.io/templ/cmds/templcommands"
	"regexp"
	"strings"
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
	tflags := flag.NewFlagSet("templFlags", flag.ContinueOnError)
	// Create a copy of os.Args
	argsCopy := append([]string(nil), os.Args...)

	templCommander := subcommands.NewCommander(tflags, os.Args[0])

	logrus.Info("Starting ", os.Args[0])

	configDir := getTemplatesDirectory()
	changeDir(configDir)

	//Super cheap way to get documentation into the usage message.
	oldHelp := templCommander.Explain
	help := func(w io.Writer) {
		fmt.Fprintf(w, "Env Vars: LOG_LEVEL TEMPL_DIR\n")
		oldHelp(w)

	}

	templateConfigPaths, argsCopy := parseArgsForDynamicConfigOptions(argsCopy)

	templCommander.Explain = help

	templCommander.Register(subcommands.HelpCommand(), "help")
	templCommander.Register(&templcommands.CatCommand{}, "templates")
	templCommander.Register(&templcommands.ListCommand{}, "templates")
	templCommander.Register(templcommands.NewRenderCommand(templateConfigPaths), "templates")

	tflags.Parse(argsCopy[1:])
	exitval := int(templCommander.Execute(ctx, os.Args))
	os.Exit(exitval)
}

// Validates that a directory exists and then changes the pwd into it.
// Parameters:
// 1. templatedir - a string representing the path to the directory containing templates.
// Returns:
// string - Absolute path to the directory that has been switched into.
// Side Effects
// Will panic if the directory does not exist.
func changeDir(templatedir string) string {
	stat, err := os.Stat(templatedir)

	if !stat.IsDir() || err != nil {
		logrus.Debug("Trying to chdir into <", templatedir, "> failed. Panicking...")
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

// Interrogates the environment variable TEMPL_DIR for a directory of templates.
// If that variable does not exist, it attempts to switch into the default template directory.
func getTemplatesDirectory() string {

	configDir, foundEnvVar := os.LookupEnv("TEMPL_DIR")

	if foundEnvVar {
		logrus.Debug("Found TEMPL_DIR=", configDir)
		return configDir
	}

	home, foundHome := os.LookupEnv("HOME")

	if !foundHome {
		fmt.Print("HOME env var not found. templ cannot autodiscover the default configdir, and TEMPL_DIR is not set.")
		panic("Please review the help for preconditions to meet running templ/")
	}

	configDir = home + "/.config/templ"

	if _, err := os.Stat(configDir); err == nil || os.IsNotExist(err) {
		logrus.Info("TEMPL_DIR not set. Did not find the default directory. Creating...")
		err := os.MkdirAll(configDir, 0700)

		if err != nil {
			panic(err)
		}
	}

	return configDir

}

// The config files for templates are a bit of a pain: they options are not registerable
// as flags, because we don't know the flag names, and the flags library has
// a strong reliance on knowing these ahead of time. When flag.Parse() is called, if a flag appears in argv that is not
// defined already then an error is raised. I tried setting the flags library to continue
// on error but this is a special case.
// Therefore I need to parse out the --config-foo options and remove them from os.args ahead of time.
// We're looking for values of type --config-<some template> path/to/config.yml
//
// The flags library allows for these parsings for args that take string vals
//
//	-flag=x
//	--flag=x
//	-flag x
//	--flag x
//
// We can imagine a template file might be named
// config-my-thing-that-takes-config-files
// so verify there's a dash at the beginning, and
// the string config followed by a dash
func parseArgsForDynamicConfigOptions(argsCopy []string) (map[string]string, []string) {
	var templateConfigPaths = make(map[string]string)

	// I can make the paths into absolute paths by prefixing them with the templates directory.
	configDir := getTemplatesDirectory()
	// remove the config file flags; these are dynamically parsed, not statically configured,
	for i, flag := range argsCopy {
		var templateName string
		var configPath string

		if strings.Contains(flag, "-config") {
			// Remove leading "--" or "-"
			re := regexp.MustCompile("^[-]+")
			result := re.ReplaceAllString(flag, "")

			templateName = strings.Split(result, "-")[0]

			if strings.Contains(flag, "=") {
				parsed := strings.Split(flag, "=")
				configPath = parsed[1]
				// Remove --foo-config=bar from args
				argsCopy = append(argsCopy[:i], argsCopy[i+1:]...)
			} else {
				configPath = os.Args[i+1]
				// Remove --foo-config bar from args
				argsCopy = append(argsCopy[:i], argsCopy[i+2:]...)
			}

			templateConfigPaths[templateName] = configDir + "/" + configPath
		}

	}
	return templateConfigPaths, argsCopy
}
