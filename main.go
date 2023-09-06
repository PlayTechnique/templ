package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"templ/configelements"
	"templ/repository"
	"templ/templatedirectories"
	"templ/templates"
)

func main() {
	// string flag called clone. Takes a url as an argument
	url := flag.String("fetch", "", "clone a git repository from a url. Can be a github url or a local git repository.")
	update := flag.Bool("update", false, "iterate over template repositories, calling git update.")

	usage := fmt.Sprintf("%s <templatename || templatename=variablesfile.yaml> <flags>\n\n"+
		"<templatename> can come in one of two forms. First is a template filename,"+
		"second is a template filename plus a yaml variables file of key: value pairs which"+
		"will populate the template file's variables. If no variables file is provided, the utility"+
		"operates like 'cat' on the file, printing it to stdout.\n\n"+
		"This utility can also be called in a pipeline as %s templatename | %s FOO=BAR BAM=BAS, for folks who"+
		"prefer not to have a variables file.\n\n",
		filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))

	flag.Usage = func() { fmt.Println(usage); flag.PrintDefaults(); return }
	flag.Parse()

	candidateTemplates := flag.Args()

	createTemplDir()

	// If the user has provided a url, clone the repo
	if *url != "" {
		// I use github exclusively right now, so this is a safe bet. If I need to support more version control systems
		// I'll need to implement a case statement and switch on some piece of data to figure out the difference between
		// repo types. Probably need an extra cli flag.
		templRepo, err := repository.NewGitRepository(*url)

		if err != nil {
			panic(err)
		}

		err = templRepo.Fetch()

		if err != nil {
			panic(err)
		}
	}

	if *update {
		err := templatedirectories.Update()

		if err != nil {
			panic(err)
		}
	}

	templates.Render(candidateTemplates)
}

//Helper functions

// createTemplDir requests information about the right path for templ's templates directory and creates that directory
// if need be. After calling this, our initial precondition should be met.
func createTemplDir() {
	templDir := configelements.NewTemplDir().TemplatesDir
	templDir, err := filepath.Abs(templDir)

	if err != nil {
		fmt.Printf("Tried to determine absolute path to the templ directory. Failed: %v", err)
		panic("")
	}

	_, err = os.Stat(templDir)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(templDir, 0700)

			if err != nil {
				fmt.Printf("Could not create configuration directory %s: %v", templDir, err)
			}
		}

		if err != nil {
			panic(err)
		}
	}
}
