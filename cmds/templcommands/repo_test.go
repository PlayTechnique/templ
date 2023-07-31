package templcommands

import (
	"testing"
)

func TestExtractDirNameFromUrl(t *testing.T) {
	// Verify the extractDirNameFromUrl function
	// it currently seems to extract both

	expected := "gwynforthewyn-jinx"
	actual, err := extractDirNameFromUrl("https://github.com/gwynforthewyn/jinx")

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if expected != actual {
		t.Errorf("Expected: %q, Got: %q", expected, actual)
	}

}
