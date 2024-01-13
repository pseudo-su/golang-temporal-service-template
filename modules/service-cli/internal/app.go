package internal

import (
	"log/slog"

	"github.com/pseudo-su/golang-temporal-service-template/modules/service-cli/internal/commands/example"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/logsetup"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:     "cli",
		Usage:    "CLI utility to interact with the golang-temporal-service-template",
		HelpName: "cli",
		Commands: []*cli.Command{
			example.Command(),
		},
		Before: func(cc *cli.Context) error {
			// Create logger
			logger := logsetup.NewLogger(
				logsetup.WithModeStr("text"),
				logsetup.WithLevelStr("info"),
			)
			slog.SetDefault(logger)

			return nil
		},
		After: func(cc *cli.Context) error {
			return nil
		},
	}
}
