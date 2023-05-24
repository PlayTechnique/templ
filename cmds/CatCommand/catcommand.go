package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
)

type CatCommand struct {
	TemplateName string
	synopsis     string
	usage        string
	mew          bool
}

var catcommand CatCommand

func init() {
	catcommand.TemplateName = "cat"
	catcommand.synopsis = "cat a named template templates"
	catcommand.usage = `
cat <name of a template file> 
Output the contents of named, known template files.

You can see the names of the template files with the list subcommand.
`
}

func (CatCommand) Name() string {
	return catcommand.TemplateName
}

func (CatCommand) Synopsis() string {
	return catcommand.synopsis
}

func (CatCommand) Usage() string {
	return catcommand.usage
}

func (c *CatCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.mew, "say a mew", false, "gotta mew before purring")
}

func (*CatCommand) Execute(_ context.Context, _ *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	for _, arg := range args {
		// search for each named template
		if arg == "kitty" {
			fmt.Println("kitty!")
		}
		// print it to stdout
	}
	return subcommands.ExitSuccess
}
