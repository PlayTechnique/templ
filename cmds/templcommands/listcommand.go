package templcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"strings"
)

type ListCommand struct {
	TemplateName string
	synopsis     string
	usage        string
}

var listCommand ListCommand

func init() {
	listCommand.TemplateName = "list"
	listCommand.synopsis = "list available templates"
	listCommand.usage = `
list 
Outputs the names of all known template files.
`
}

func (ListCommand) Name() string {
	return listCommand.TemplateName
}

func (ListCommand) Synopsis() string {
	return listCommand.synopsis
}

func (ListCommand) Usage() string {
	return listCommand.usage
}

func (c *ListCommand) SetFlags(_ *flag.FlagSet) {

}

func (*ListCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logrus.Debug(f)
	files, err := listFiles([]string{"."})

	fmt.Println(strings.Join(files, "\n"))

	return err
}
