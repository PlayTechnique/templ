package configelements_test

import (
	"os"
	"templ/configelements"
	"testing"
)

func TestGettingDefaultTemplDir(t *testing.T) {
	templDir := configelements.NewTemplDir().TemplatesDir
	defaultDir := os.Getenv("HOME") + "/.config/templ"

	if templDir != defaultDir {
		t.Errorf("Calculated default dir <%s> does not equal expected default <%s>", templDir, defaultDir)
	}
}

func TestSettingTemplDirViaEnv(t *testing.T) {
	//The implementation of TemplatesDir only returns absolute paths, so this has to be absolute so the comparison works.
	//Note that at this point the template dir does not have to exist.
	newTemplDir, err := os.MkdirTemp("", "templ-testing-setting-via-env-var")

	if err != nil {
		t.Errorf("TestSettingTemplDirViaEnv could not create templates dir")
	}

	err = os.Setenv("TEMPL_DIR", newTemplDir)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := os.Unsetenv("TEMPL_DIR"); err != nil {
			t.Errorf("Error unsetting %s env variable: %v\n", "TEMPL_DIR", err)
		}

		if err = os.Remove(newTemplDir); err != nil {
			t.Errorf("Error deleting temp directory %s: %v\n", newTemplDir, err)
		}
	}()

	templDir := configelements.NewTemplDir().TemplatesDir

	if templDir != newTemplDir {
		t.Errorf("Templ dir is <%s>, should be <%s>", templDir, newTemplDir)
	}
}
