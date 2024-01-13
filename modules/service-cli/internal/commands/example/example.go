package example

import (
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "example",
		Usage:       "Subcommand example",
		Description: "Subcommand example",
		Subcommands: []*cli.Command{
			printHelloWorldCommand(),
		},
	}
}
