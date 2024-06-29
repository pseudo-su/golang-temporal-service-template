package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/initialise"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/internal"
)

func main() {
	ctx := context.Background()

	cfg, err := initialise.ServiceWithConfig(ctx, &internal.WorkerConfig{})
	if err != nil {
		slog.ErrorContext(context.Background(), "unable to parse environment variables", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "App config loaded", slog.Any("name", cfg.App.Name), slog.Any("env", cfg.App.Env))

	w, err := internal.NewWorker(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "unable to create worker", slog.Any("error", err))
		os.Exit(1)
	}
	w.Run(ctx)

	slog.InfoContext(ctx, "clean shutdown")
	os.Exit(0)
}
