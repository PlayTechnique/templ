package test_helpers

import (
	"os"
	"testing"
)

func CleanUpTemplDir(tempDir string, t *testing.T) {
	err := os.Unsetenv("TEMPL_DIR")

	if err != nil {
		t.Errorf("could not unset TEMPL_DIR env var: %v", err)
	}

	err = os.RemoveAll(tempDir)

	if err != nil {
		t.Errorf("could not remove temp dir: %v", err)
	}
}
