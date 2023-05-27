package templcommands

import "os"

type TestFileStructure struct {
	directories []string
	files       []string
}

type TestSetup struct {
	name       string
	setupFiles TestFileStructure
	startDirs  []string
	want       []string
}

func (t TestSetup) Setup() (tempdir string) {
	tempdir, err := os.MkdirTemp("", "templ_test")

	if err != nil {
		panic(err)
	}

	err = os.Chdir(tempdir)
	if err != nil {
		panic(err)
	}

	for _, dir := range t.setupFiles.directories {
		os.MkdirAll(dir, 0755)
	}

	for _, files := range t.setupFiles.files {
		os.Create(files)
	}

	return tempdir
}

func (t TestSetup) TearDown(directory string) {
	err := os.RemoveAll(directory)
	if err != nil {
		panic(err)
	}

}
