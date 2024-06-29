package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/grpchealth"
	"github.com/pseudo-su/golang-temporal-service-template/modules/mocks/internal"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/httpserver"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/initialise"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/rungroup"
)

func main() {
	ctx := context.Background()

	cfg, err := initialise.ServiceWithConfig(ctx, &internal.MocksConfig{})
	if err != nil {
		slog.ErrorContext(context.Background(), "unable to parse environment variables", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "App config loaded", slog.Any("name", cfg.App.Name), slog.Any("env", cfg.App.Env))

	slog.InfoContext(ctx, "Initialising http server")
	address := fmt.Sprintf(":%d", cfg.Tcp.Port)
	httpServer, err := httpserver.New(
		ctx,
		httpserver.WithAddress(address),
		httpserver.WithReflection(),
		// httpserver.WithUnaryInterceptors(),
		// httpserver.WithStreamInterceptors(),
	)
	if err != nil {
		slog.ErrorContext(ctx, "error creating server", slog.Any("error", err))
		os.Exit(1)
	}

	httpServer.RegisterConnectHandler(func(connectServer *http.ServeMux) {
		connectServer.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	})

	// Run server
	rg := rungroup.NewRunGroup(ctx)
	rg.Run(func(ctx context.Context) {
		slog.InfoContext(ctx, "starting http server")
		if err := httpServer.ListenAndServe(rungroup.InterruptChannel(ctx)); err != nil {
			slog.ErrorContext(ctx, "http server error", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "http server stopped")
	})
	rg.Wait()
	slog.InfoContext(ctx, "clean shutdown")
	os.Exit(0)
}
