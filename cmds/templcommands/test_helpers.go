package templcommands

import (
	"os"
)

type TestFileStructure struct {
	directories []string
	files       map[string]string
}

type TestSetup struct {
	name       string
	setupFiles TestFileStructure
	startDirs  []string
	want       []string
}

func Setup(t TestFileStructure) (tempdir string) {
	tempdir, err := os.MkdirTemp("", "templ_test")

	if err != nil {
		panic(err)
	}

	err = os.Chdir(tempdir)

	if err != nil {
		panic(err)
	}

	for _, dir := range t.directories {
		os.MkdirAll(dir, 0755)
	}

	for filename, content := range t.files {

		err := os.WriteFile(filename, []byte(content), 0755)

		if err != nil {
			panic(err)
		}
	}

	return tempdir
}

func TearDown(directory string) {
	err := os.RemoveAll(directory)
	if err != nil {
		panic(err)
	}

}
