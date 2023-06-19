package templcommands

import (
	"context"
	"errors"
	_ "errors"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	_ "os"
	"strings"
	"text/template"
)

type templatePath string
type templateVariablesPath string

type RenderCommand struct {
	templateName        string
	templatedirectory   string
	templateconfigpaths map[templatePath]templateVariablesPath
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
render <name of at least one template file from the list>=<path to template variables file>

Output the contents of named, known template files. You supply the variables through a yaml file. It should contain a map.

If we have a template file called examples/hello.txt:
=======
Hello {{ .name }}
=======

Then the yaml file, in path /path/to/hello.yaml should read:
=======
---
name: "World"
=======

Then you run:
render examples/hello.txt=/path/to/hello.yaml

And the output is:
Hello World


To find out the names of available templates, use the list subcommand.

All variables in the template must be populated in the template config file. To turn this behaviour off, use the --no-strict option.

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

func (r RenderCommand) Execute(c context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 0 {
		logrus.Error("No template files given")
		fmt.Print(r.Usage())
		return subcommands.ExitFailure
	}

	// The version of argv returned by flags.args should only hold arguments now, which is the template files
	var templateFiles = make([]templatePath, len(f.Args()))
	var templateVariablesFiles = make(map[templatePath]templateVariablesPath, len(f.Args()))

	for _, combo := range f.Args() {
		logrus.Debug("combo: ", combo)
		//split the combo into template file and template variables file around the = sign
		split := strings.Split(combo, "=")

		if len(split) != 2 {
			logrus.Error("Failed to split template file and template variables file. Input was: ", combo)
			return subcommands.ExitFailure
		}

		var tp = templatePath(split[0])
		templateFiles = append(templateFiles, tp)
		templateVariablesFiles[tp] = templateVariablesPath(split[1])
	}

	exitStatus, _ := render(templateFiles, rendercommand.templateconfigpaths)
	return exitStatus
}

func render(templateFiles []templatePath, templateVariables map[templatePath]templateVariablesPath) (subcommands.ExitStatus, error) {
	err := validateTemplatesExist(templateFiles)

	if err != nil {
		logrus.Error(err)
		return subcommands.ExitFailure, err
	}

	logrus.Debug("filesInArgs: ", templateFiles)

	for _, templatePath := range templateFiles {
		templateVariablesFilePath := templateVariables[templatePath]

		if templateVariablesFilePath == "" {
			logrus.Debug("No template variables file given for template file: ", templatePath)
			continue
		} else {
			logrus.Debug("Found template variables file: ", templateVariablesFilePath)
		}
		// Consume the template variables, which are a yaml file, into a map
		// of key value pairs.
		templateVariables, err := getTemplateVariables(templateVariablesFilePath)

		// Read the template file
		templateContents, err := os.ReadFile(string(templatePath))

		if err != nil {
			logrus.Info("Failed to read template file: ", templatePath, " with error ", err)
			return subcommands.ExitFailure, err
		} else {
			logrus.Info("Successfully read template file: ", templatePath)
		}

		// Convert template file content to a string
		templateText := string(templateContents)

		// Create a new template and parse the template text
		tmpl, err := template.New(string(templatePath)).Parse(templateText)

		if err != nil {
			logrus.Error("Failed to parse template: ", err)
			return subcommands.ExitFailure, err
		} else {
			logrus.Info("Successfully parsed template: ", templatePath)
		}

		// Execute the template with the data
		err = tmpl.Execute(os.Stdout, templateVariables)

		if err != nil {
			log.Fatalf("Failed to execute template: %v", err)
			return subcommands.ExitFailure, err
		} else {
			logrus.Info("Successfully executed template: ", templatePath)
		}
	}

	return subcommands.ExitSuccess, nil
}

func getTemplateVariables(templateVariablesFilePath templateVariablesPath) (map[string]string, error) {
	// Read the YAML file

	yamlFile, err := os.ReadFile(string(templateVariablesFilePath))
	if err != nil {
		logrus.Error("Failed to read YAML file at path <", templateVariablesFilePath, "> Err is: ", err)
		return nil, err
	}

	// Create a map to store the parsed YAML data
	data := make(map[string]string)

	// Unmarshal the YAML data into the map
	err = yaml.Unmarshal(yamlFile, &data)

	if err != nil {
		logrus.Error("Failed to unmarshal YAML: ", err)
		return nil, err
	}

	return data, err
}

func validateTemplatesExist(templateFiles []templatePath) error {

	// First valid named files in the arguments.
	// Then re-scan across os.args to find any other flags
	// that contain the filenames of the files.
	for _, templateFilePath := range templateFiles {
		logrus.Debug(" variable:", string(templateFilePath))
		file, err := os.OpenFile(string(templateFilePath), os.O_RDONLY, 0644)
		defer file.Close()

		if errors.Is(errors.Unwrap(err), os.ErrNotExist) {
			// The argument is not a file. Proceed to
			// handle the case where the file doesn't exist
			err = fmt.Errorf("File %s does not exist: %v", templateFilePath, err)
			return err
		}

	}

	return nil
}

func NewRenderCommand(templateDir string) RenderCommand {

	rendercommand.templatedirectory = templateDir
	return rendercommand
}
