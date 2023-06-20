package templcommands

import (
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func init() {
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

func TestRenderFile(t *testing.T) {
	// Set up a test matrix
	testCases := []TestSetup{
		{name: "No argument given", setupFiles: TestFileStructure{directories: []string{}, files: map[string]string{}}, startDirs: []string{"./"}, want: []string{}},
	}

	for _, tt := range testCases {

		tempdir := Setup(tt.setupFiles)

		err := os.Chdir(tempdir)
		if err != nil {
			panic(err)
		}

		defer TearDown(tempdir)

		t.Run(tt.name, func(t *testing.T) {
			ans, err := listFiles(tt.startDirs)

			assert.True(t, err == subcommands.ExitSuccess, "Expected success, got %s", err)

			sort.Strings(ans)
			sort.Strings(tt.want)
			assert.True(t, reflect.DeepEqual(ans, tt.want), "Expected %s, got %s", tt.want, ans)

		},
		)
	}
}

func TestRender(t *testing.T) {
	setupStructure := TestFileStructure{
		directories: []string{"level_one", "level_one/level_two"},
		files: map[string]string{
			"ringding":              "this is the content",
			"level_one/smudge":      "more content {{ .roflcopter }} here\n",
			"level_one/smudge.yaml": "---\nroflcopter: \"hippololamus\"\n",
		},
	}

	tempdir, templatePaths, templateVariables := setupRenderTestStructures(setupStructure)

	err := os.Chdir(tempdir)

	if err != nil {
		panic(err)
	}
	defer TearDown(tempdir)

	exitstatus, _ := render(templatePaths, templateVariables)
	assert.IsTypef(t, subcommands.ExitSuccess, exitstatus, "Expected success, got %s", exitstatus)
}

func setupRenderTestStructures(setupStructure TestFileStructure) (string, []templatePath, map[templatePath]templateVariablesPath) {

	tempdir := Setup(setupStructure)
	err := os.Chdir(tempdir)
	if err != nil {
		panic(err)
	}

	// The render function accepts a structure of type templatePath and a structure of type templateVariablesPath.
	var templatePaths = make([]templatePath, 0)
	var templateVariables = make(map[templatePath]templateVariablesPath, 0)

	for filename, fileContent := range setupStructure.files {

		err := os.WriteFile(tempdir+"/"+filename, []byte(fileContent), 0644)

		if err != nil {
			panic(err)
		}

		// If this is a config file, then we store its path in the templateVariables data structure
		if isVariablesFile(filename) {
			associatedTemplateFile := templatePath(strings.TrimSuffix(filename, ".yaml"))
			templateVariables[associatedTemplateFile] = templateVariablesPath(filename)
		}

		// If this is a template file, the templatePaths data structure is index based
		if !isVariablesFile(filename) {
			templatePaths = append(templatePaths, templatePath(filename))
		}

	}

	return tempdir, templatePaths, templateVariables
}

func TestRenderFileThatDoesNotExist(t *testing.T) {
	exit, err := render([]templatePath{templatePath("does_not_exist")}, map[templatePath]templateVariablesPath{"does_not_exist": "does_not_exist"})
	assert.IsTypef(t, subcommands.ExitFailure, exit, "Expected failure, got %s", err)
	assert.Errorf(t, err, "open does_not_exist: no such file or directory")
}

func isVariablesFile(filename string) bool {
	return strings.HasSuffix(filename, ".yaml")
}
