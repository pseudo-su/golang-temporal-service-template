package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/pseudo-su/golang-temporal-service-template/modules/service-cli/internal"
)

func main() {
	ctx := context.Background()
	app := internal.NewApp()

	err := app.Run(os.Args)
	if err != nil {
		slog.ErrorContext(ctx, "cli exit 1", slog.Any("error", err))
		os.Exit(1)
	}
}
