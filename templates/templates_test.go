package templates_test

import (
	"templ/templates"
	"testing"
)

func TestRenderFromStringWithEmptyString(t *testing.T) {
	hydratedTemplate, err := templates.RenderFromString("", []string{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != "" {
		t.Errorf("HydratedTemplate should be empty. Received <%s>", hydratedTemplate)
	}
}

func TestRenderFromStringWithPlainString(t *testing.T) {
	template := `
I love humans
`
	hydratedTemplate, err := templates.RenderFromString(template, []string{})

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != template {
		t.Errorf("HydratedTemplate should be contain a simple string. Received <%s>", hydratedTemplate)
	}
}

func TestRenderFromStringWithTemplateContainingAVariableButNoVariables(t *testing.T) {
	template := `I love {{ .SPECIES }}`

	templateVariables := []string{}
	hydratedTemplate, err := templates.RenderFromString(template, templateVariables)

	if err != nil {
		t.Errorf("%v", err)
	}

	if hydratedTemplate != template {
		t.Errorf("Expected <%s>, received <%s>", template, hydratedTemplate)
	}
}

func TestRenderFromStringWithVariablesAndDefinitions(t *testing.T) {
	template := `I love {{ .SPECIES }}`

	templateVariables := []string{"SPECIES=HUMAN"}
	hydratedTemplate, err := templates.RenderFromString(template, templateVariables)

	if err != nil {
		t.Errorf("%v", err)
	}

	expected := "I love HUMAN"
	if hydratedTemplate != expected {
		t.Errorf("Expected <%s>, received <%s>", expected, hydratedTemplate)
	}
}
