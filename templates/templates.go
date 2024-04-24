package templates

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"templ/configelements"
	"text/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type TemplateVariableErr struct {
	ErrorMessage string
}

func (t TemplateVariableErr) Error() string {
	return t.ErrorMessage
}

func (t TemplateVariableErr) Is(target error) bool {
	_, ok := target.(TemplateVariableErr)
	return ok
}

func RenderFromStdin(template string, variableDefinitions []string) (hydratedtemplate string, err error) {
	// If we don't receive any arguments, just pass back up the chain.
	if len(variableDefinitions) == 0 {
		return template, nil
	}

	variables, err := convertFromArrayToKeymap(variableDefinitions)

	if err != nil {
		return "", err
	}

	hydratedtemplate, err = renderFromString("stdin", template, variables)

	return
}

func FindTemplateAndVariableFiles(argv []string) ([]string, map[string]string, error) {
	// Data structures to store paths to the template files. These may optionally have an associated variables file to hydrate with.
	var templateFilePaths = make([]string, 0)
	var templateVariablesFilesPaths = make(map[string]string, 0)

	for _, template := range argv {
		// The arguments at this point either read as a name/of/template_file, or as name/of/template_file=path/to/variables.
		// In the first case, I want to store the path to the template file in an array to hand in to the renderFromFiles command.
		// In the second case, we store the path to the template file in the same array, and also use that path as an
		// index in a map, where the value is the variables file's path.

		// interrogate each path from args to split into either a set of strings or a path+string
		// then use findFilesByName to find the templates associated with those strings
		variablesFile := false
		var templateVariablesPath string

		if strings.Contains(template, "=") {
			variablesFile = true
			templateAndVariablesPath := strings.Split(template, "=")
			template = templateAndVariablesPath[0]
			templateVariablesPath = templateAndVariablesPath[1]
		}

		//temp variable to prevent variable shadowing.
		t, err := findFilesByName(configelements.NewTemplDir().TemplatesDir, []string{template})

		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		} else {
			logrus.Debug("Found templateFilePaths,", templateFilePaths)
		}

		for _, tp := range t {
			p := tp

			templateFilePaths = append(templateFilePaths, p)

			if variablesFile {
				templateVariablesFilesPaths[p] = templateVariablesPath
			}
		}
	}
	return templateFilePaths, templateVariablesFilesPaths, nil
}

func RenderFromFiles(templateFiles []string, templateVariables map[string]string) error {
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
				_, file, line, _ := runtime.Caller(0)
				return fmt.Errorf("%s:%d: %v", file, line, err)
			}

			fmt.Println(string(content))

			continue
		}

		logrus.Debug("Found template variables file: ", templateVariablesFilePath)

		// Consume the template variables, which are a yaml file, into a map
		// of key value pairs.
		templateVariables, err := getTemplateVariablesFromYamlFile(templateVariablesFilePath)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("%s:%d: %v", file, line, err)
		}

		templateContents, err := os.ReadFile(templatePath)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("%s:%d: %v", file, line, err)
		}

		// Convert template file content to a string
		templateText := string(templateContents)
		output, err := renderFromString(templatePath, templateText, templateVariables)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("%s:%d: %v", file, line, err)
		}

		fmt.Println(output)
	}

	return nil
}

// RetrieveVariables accepts the content of a template and returns an array of.
// strings that match {{ FOO }} format, but not with formats like ${{ FOO }}
func RetrieveVariables(templateContent string) []string {
	// Regular expression for strict matches: only `{{ .Identifier }}`
	strictRe := regexp.MustCompile(`{{\s*\.([^}\s]+)\s*}}`)
	strictMatches := strictRe.FindAllStringSubmatch(templateContent, -1)

	matches := []string{}

	for _, section := range strictMatches {
		if !slices.Contains(matches, section[1]) {
			matches = append(matches, section[1])
		}
	}

	return matches
}

// renderFromString takes a string containing a template, and a map of FOO=bar definitions and returns
// a rendered template.
// Problematic workflow:
// Templates are arbitrary files, and a lot of template formats use some variation on the {{}} delimiters to identify
// strings that should be substituted. When parsing these, the golang template engine sees these double braces and tries
// to replace these itself, which causes errors.
// For example, here's a subsection of a github workflow:
// step: "do something cool"
// =====
// run : |
//
//	sed {{ .FILENAME }} -i /${{ env.PATTERN }}/sub'
//
// =====
// In this file, I want templ to edit {{ .FILENAME }}, but the section ${{ env.PATTERN }} will cause an error because templ
// isn't handing in an object called env with a PATTERN member.
// Solution:
// To work around this, templ splits all incoming files on the string '{{'. It then either successfully substitutes
// a variable or it ignores any error rendering that subsection. Tada!
func renderFromString(templatePath string, templateText string, templateVariableDefinitions map[string]string) (string, error) {

	templateSections := strings.SplitAfter(templateText, "}}")
	var reformedTemplate bytes.Buffer

	templateName := path.Base(templatePath)

	for i, section := range templateSections {
		name := templateName + strconv.Itoa(i)

		tmpl, err := template.New(name).Parse(section)

		if err != nil {
			// If the error reads like this:
			// `variable env.FOO not defined`
			// then we're in the error condition identified in the header comment. Ignore the error, write the contents
			// of the template section to our template buffer and go on to the next loop.
			if strings.Contains(err.Error(), "not defined") {
				reformedTemplate.WriteString(section)
				continue
			} else if strings.Contains(err.Error(), "bad character") {
				_, file, line, _ := runtime.Caller(0)
				return "", fmt.Errorf("\nThe 'bad character' error from the go template engine normally means that you have a disallowed "+
					"character inside your template variable name. \nHere's what templ was working on:\n"+
					"%s:%d: Failed on section:\n--- %s\n---\n Error is: %v", file, line, section, err)
			} else {
				_, file, line, _ := runtime.Caller(0)
				return "", fmt.Errorf("%s:%d: Failed on section:\n--- %s\n---\n Error is: %v", file, line, section, err)

			}
		}

		// Execute the template with the data
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, templateVariableDefinitions)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return buffer.String(), fmt.Errorf("%s:%d: %v", file, line, err)
		}

		reformedTemplate.Write(buffer.Bytes())
	}

	return reformedTemplate.String(), nil
}

func getTemplateVariablesFromYamlFile(templateVariablesFilePath string) (map[string]string, error) {
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
		// The argument is not a file, which can be normal, but if
		// raise if the file doesn't exist
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("file %s does not exist: %v", templateFilePath, err)
			return err
		}

		defer file.Close()

	}

	return nil
}

// TODO: Why do template rendering functions care where the files are? Move this to main or somewhere.
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

func convertFromArrayToKeymap(input []string) (map[string]string, error) {
	k := make(map[string]string)

	for _, arg := range input {
		if !strings.Contains(arg, "=") {
			_, file, line, _ := runtime.Caller(0)

			message := fmt.Sprintf("%s:%d: Argument <%s> not formatted as FOO=BAR", file, line, arg)
			err := TemplateVariableErr{ErrorMessage: message}
			return k, err
		}

		s := strings.Split(arg, "=")
		k[s[0]] = s[1]
	}

	return k, nil
}
