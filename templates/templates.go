package templates

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"templ/configelements"
	"text/template"
)

func Render(templates []string) error {
	// Data structures to store paths to the template files. These may optionally have an associated variables file to hydrate with.
	var templateFilePaths = make([]string, 0)
	var templateVariablesFilesPaths = make(map[string]string, 0)

	for _, path := range templates {
		// The arguments at this point either read as a name/of/template_file, or as name/of/template_file=path/to/variables.
		// In the first case, I want to store the path to the template file in an array to hand in to the render command.
		// In the second case, we store the path to the template file in the same array, and also use that path as an
		// index in a map, where the value is the variables file's path.

		// interrogate each path from args to split into either a set of strings or a path+string
		// then use findFilesByName to find the templates associated with those strings
		variablesFile := false
		var templateVariablesFilePath string

		if strings.Contains(path, "=") {
			variablesFile = true
			templatePathAndVariablesPath := strings.Split(path, "=")
			path = templatePathAndVariablesPath[0]
			templateVariablesFilePath = templatePathAndVariablesPath[1]
		}

		//temp variable to prevent variable shadowing.
		t, err := findFilesByName(configelements.NewTemplDir().TemplatesDir, []string{path})

		if err != nil {
			logrus.Error(err)
			return err
		} else {
			logrus.Debug("Found templateFilePaths,", templateFilePaths)
		}

		for _, tp := range t {
			p := tp

			templateFilePaths = append(templateFilePaths, p)

			if variablesFile {
				templateVariablesFilesPaths[p] = templateVariablesFilePath
			}
		}
	}

	err := render(templateFilePaths, templateVariablesFilesPaths)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func render(templateFiles []string, templateVariables map[string]string) error {
	err := validateTemplatesExist(templateFiles)

	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Debug("filesInArgs: ", templateFiles)

	for _, templatePath := range templateFiles {
		templateVariablesFilePath := templateVariables[templatePath]

		// No variables? Just print and move on.
		if templateVariablesFilePath == "" {
			content, err := os.ReadFile(string(templatePath))

			if err != nil {
				logrus.Error("Failed to read template file: ", templatePath, " with error ", err)
				continue
			}

			fmt.Println(string(content))

			continue
		}

		logrus.Debug("Found template variables file: ", templateVariablesFilePath)

		// Consume the template variables, which are a yaml file, into a map
		// of key value pairs.
		templateVariables, err := getTemplateVariables(templateVariablesFilePath)

		templateContents, err := os.ReadFile(string(templatePath))

		if err != nil {
			logrus.Info("Failed to read template file: ", templatePath, " with error ", err)
			return err
		} else {
			logrus.Info("Successfully read template file: ", templatePath)
		}

		// Convert template file content to a string
		templateText := string(templateContents)

		// Create a new template and parse the template text
		tmpl, err := template.New(string(templatePath)).Parse(templateText)

		if err != nil {
			logrus.Error("Failed to parse template: ", err)
			return err
		} else {
			logrus.Info("Successfully parsed template: ", templatePath)
		}

		// Execute the template with the data
		err = tmpl.Execute(os.Stdout, templateVariables)

		if err != nil {
			log.Fatalf("Failed to execute template: %v", err)
			return err
		} else {
			logrus.Info("Successfully executed template: ", templatePath)
		}
	}

	return nil
}

func getTemplateVariables(templateVariablesFilePath string) (map[string]string, error) {
	// Read the YAML file
	yamlFile, err := os.ReadFile(templateVariablesFilePath)
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

func validateTemplatesExist(templateFiles []string) error {

	// First valid named files in the arguments.
	// Then re-scan across os.args to find any other flags
	// that contain the filenames of the files.
	for _, templateFilePath := range templateFiles {
		logrus.Debug(" variable:", string(templateFilePath))
		file, err := os.OpenFile(string(templateFilePath), os.O_RDONLY, 0644)
		defer file.Close()

		if errors.Is(err, os.ErrNotExist) {
			// The argument is not a file. Proceed to
			// handle the case where the file doesn't exist
			err = fmt.Errorf("File %s does not exist: %v", templateFilePath, err)
			return err
		}

	}

	return nil
}

// findFilesByName searches a directory for file names that match those provided in a set of strings.
// Arguments:
// root: the directory to search
// names: a set of filenames to search for
// Returns:
// an array of strings, each of which is the path to a file that was found.
// or an error
func findFilesByName(root string, names []string) ([]string, error) {
	foundFiles := []string{}

	logrus.Debug("Outside filepath.Walk function names: ", names)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logrus.Debug("Inside filepath.Walk function names: ", names)

		// If the file's name is in the set of names
		for _, name := range names {
			logrus.Debug("name is <", name, "> path is <", path, ">")
			if strings.Contains(path, name) {
				foundFiles = append(foundFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return foundFiles, nil
}
