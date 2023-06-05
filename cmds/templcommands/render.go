package templcommands

import (
	"context"
	"errors"
	_ "errors"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	_ "os"
)

type RenderCommand struct {
	dynamicflags        []string
	templateName        string
	templateconfigpaths map[string]string
	synopsis            string
	usage               string
	strict              bool
}

//templ render flobble --flobble-config=some/path
//templ render flobble --flobble-container=roflcopter --flobble-snibble=othervalue
//--strict: fail if any of the variables are not set

var rendercommand RenderCommand

func init() {
	rendercommand.templateName = "render"
	rendercommand.synopsis = "render a template"
	rendercommand.usage = `
render [--<template-name>-config=<value> --no-strict] <name of a template file>
Output the contents of named, known template files.

All variables in the template must be populated by default. To turn this behaviour off, use the --no-strict option.

If the name of the template  is "roflcopter", then the option --roflcopter-config=<path to config file> will be used
to specify the path to the config file for the template.
`

}

func (RenderCommand) Name() string {
	return rendercommand.templateName
}

func (RenderCommand) Synopsis() string {
	return rendercommand.synopsis
}

func (RenderCommand) Usage() string {
	return rendercommand.usage
}

func (r RenderCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&r.strict, "no-strict", true, "capitalize output")
}

func (RenderCommand) Execute(c context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return render(f)
}

func render(f *flag.FlagSet) subcommands.ExitStatus {
	filesInArgs, err := findfilesFromFlags(f)
	if errors.Is(errors.Unwrap(err), os.ErrNotExist) {
		logrus.Error(err)
		return subcommands.ExitFailure
	}
	logrus.Debug("filesInArgs: ", filesInArgs)

	for _, templateFile := range filesInArgs {
		// Can I create this as a flag on the fly and get its value?
		associatedConfigFlag := templateFile + "-config"
		var r string
		f.StringVar(&r, associatedConfigFlag, "", "not sure")
		flag.Parse()

		return subcommands.ExitSuccess

	}

	return subcommands.ExitSuccess
}

func findfilesFromFlags(f *flag.FlagSet) ([]string, error) {
	filesInTheArgs := []string{}

	// First valid named files in the arguments.
	// Then re-scan across os.args to find any other flags
	// that contain the filenames of the files.
	for i, templatePath := range f.Args() {
		logrus.Debug("index: ", i, " variable:", templatePath)
		file, err := os.OpenFile(templatePath, os.O_RDONLY, 0644)
		defer file.Close()

		if errors.Is(err, os.ErrNotExist) {
			// The argument is not a file. Proceed to
			// handle the case where the file doesn't exist
			return nil, fmt.Errorf("File does not exist: %s, %w", templatePath, err)
		} else {
			filesInTheArgs = append(filesInTheArgs, templatePath)
		}
	}

	return filesInTheArgs, nil
}

func NewRenderCommand(templateConfigPaths map[string]string) subcommands.Command {
	return RenderCommand{templateconfigpaths: templateConfigPaths}
}
