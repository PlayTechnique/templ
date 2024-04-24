package templates_test

import (
	"errors"
	"slices"
	"templ/templates"
	"testing"
)

func TestTemplateVariableErr(t *testing.T) {
	e := templates.TemplateVariableErr{ErrorMessage: "test string"}
	if e.Error() != "test string" {
		t.Errorf("TemplateVariableErr not populating correctly. Expected 'test string', got %s", e.Error())
	}

	if !errors.Is(e, templates.TemplateVariableErr{}) {
		t.Errorf("TemplateVariableErr not responding as the correct type.")
	}
}

func TestRenderFromStdinWithEmptyString(t *testing.T) {
	hydratedTemplate, err := templates.RenderFromStdin("", []string{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != "" {
		t.Errorf("HydratedTemplate should be empty. Received <%s>", hydratedTemplate)
	}
}

func TestRenderFromStdingWithPlainString(t *testing.T) {
	template := `
I love humans
`
	hydratedTemplate, err := templates.RenderFromStdin(template, []string{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != template {
		t.Errorf("HydratedTemplate should be contain a simple string. Received <%s>", hydratedTemplate)
	}
}

func TestRenderFromStdinWithTemplateContainingAVariableButNoVariables(t *testing.T) {
	template := `I love {{ .SPECIES }}`

	templateVariables := []string{}
	hydratedTemplate, err := templates.RenderFromStdin(template, templateVariables)

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != template {
		t.Errorf("Expected <%s>, received <%s>", template, hydratedTemplate)
	}
}

func TestRenderFromStdinWithVariablesAndDefinitions(t *testing.T) {
	template := `I love {{ .SPECIES }}`

	templateVariables := []string{"SPECIES=HUMAN"}
	hydratedTemplate, err := templates.RenderFromStdin(template, templateVariables)

	if err != nil {
		t.Errorf("%v", err)
	}

	expected := "I love HUMAN"
	if hydratedTemplate != expected {
		t.Errorf("Expected <%s>, received <%s>", expected, hydratedTemplate)
	}
}

func TestRenderFromStdinWithVariablesAndInvalidDefinitions(t *testing.T) {
	template := `I love {{ .SPECIES }}`

	templateVariables := []string{"SPECIES HUMAN"}
	_, err := templates.RenderFromStdin(template, templateVariables)

	if !errors.Is(err, templates.TemplateVariableErr{}) {
		t.Errorf("Using an invalid variables string should raise a templatevariableerror, actually got %v", err)
	}
}

func TestRenderTemplateContainingDoubleBracesThatAreNotGoTemplateBraces(t *testing.T) {
	template := `
jobs:
  build-and-release-tag:
    env:
      OUTPUT_BINARY: {{ .BINARY_NAME }}
    steps:
      - name: "checkout"
        uses: actions/checkout@v3
        with:
          ref: ${{ env.GITHUB_REF }}
`
	templateVariables := []string{"BINARY_NAME=ROFLCOPTER"}
	hydratedTemplate, err := templates.RenderFromStdin(template, templateVariables)
	if err != nil {
		t.Errorf("%v", err)
	}
	expected := `
jobs:
  build-and-release-tag:
    env:
      OUTPUT_BINARY: ROFLCOPTER
    steps:
      - name: "checkout"
        uses: actions/checkout@v3
        with:
          ref: ${{ env.GITHUB_REF }}
`
	if hydratedTemplate != expected {
		t.Errorf("Expected <%s>, received <%s>", expected, hydratedTemplate)
	}
}

func TestCanParseTemplateVariablesFromATemplateWithoutInvalidTemplateVariables(t *testing.T) {
	template := `
jobs:
  build-and-release-tag:
    env:
      OUTPUT_BINARY: {{ .BINARY_NAME }}
    steps:
      - name: "checkout"
        uses: actions/checkout@v3
        with:
          ref: ${{ env.GITHUB_REF }}
`
	templateVariables := templates.RetrieveVariables(template)
	expected := []string{"BINARY_NAME"}

	if !slices.Equal(templateVariables, expected) {
		t.Errorf("Expected <%s>, received <%s>", expected, templateVariables)
	}
}
