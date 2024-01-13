package example

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func printHelloWorldCommand() *cli.Command {
	return &cli.Command{
		Name:   "print-hello-world",
		Action: printHelloWorldAction,
		Flags:  []cli.Flag{},
	}
}

func printHelloWorldAction(cc *cli.Context) error {
	fmt.Fprintln(os.Stdout, "Hello world")
	return nil
}
