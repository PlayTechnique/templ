package main

import (
	"flag"
	"templ/repository"
)

func main() {
	// string flag called clone. Takes a url as an argument
	url := flag.String("clone", "", "clone a repo from a url. Can be a github url or a local git repository.")

	flag.Parse()

	// I use github exclusively right now, so this is a safe bet. If I need to support more version control systems
	// I'll need to implement a case statement and switch on some piece of data to figure out the difference between
	// repo types. Probably need an extra cli flag.

	// If the user has provided a url, clone the repo
	if *url != "" {
		templRepo, err := repository.NewGitRepository(*url)

		if err != nil {
			panic(err)
		}

		err = templRepo.Fetch()

		if err != nil {
			panic(err)
		}
	}
}
