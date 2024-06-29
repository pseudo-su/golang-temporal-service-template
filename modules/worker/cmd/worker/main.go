package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/grpchealth"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/httpserver"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/initialise"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/rungroup"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/internal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	ctx := context.Background()

	cfg, err := initialise.ServiceWithConfig(ctx, &internal.WorkerConfig{})
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
		httpserver.WithUnaryInterceptors(),
		httpserver.WithStreamInterceptors(),
	)
	if err != nil {
		slog.ErrorContext(ctx, "error creating server", slog.Any("error", err))
		os.Exit(1)
	}

	// healthServer := health.NewServer()
	// httpServer.RegisterGrpcServer(func(s *grpc.Server) {
	// 	grpc_health_v1.RegisterHealthServer(s, healthServer)
	// })

	httpServer.RegisterConnectHandler(func(connectServer *http.ServeMux) {
		connectServer.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	})

	// load secrets from aegis - not yet

	slog.InfoContext(ctx, "Initialising temporal client")
	tq := cfg.Temporal.TaskQueue
	tc, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Uri.AsString(),
		Namespace: cfg.Temporal.Namespace,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error creating temporal client", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Initialising temporal worker")
	w := worker.New(tc, tq, worker.Options{})

	// register workflows - not yet

	// register activities - not yet

	// Run server
	rg := rungroup.NewRunGroup(ctx)
	rg.Run(func(ctx context.Context) {
		slog.InfoContext(ctx, "starting worker")
		if err = w.Run(rungroup.InterruptChannel(ctx)); err != nil {
			slog.ErrorContext(ctx, "worker error", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "worker stopped")
	})
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
