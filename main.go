package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
	templcommands "playtechnique.io/templ/cmds/CatCommand"
)

type Inventory struct {
	Material string
	Count    uint
}

func main() {
	ctx := context.Background()
	subcommands.Register(subcommands.HelpCommand(), "help")
	subcommands.Register(&templcommands.CatCommand{}, "templates")

	// Mystical. This seems to parse the subcommand flags.
	flag.Parse()
	os.Exit(int(subcommands.Execute(ctx, os.Args[2:])))
}
