package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/grpchealth"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/httpserver"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/rungroup"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/internal/workflows/deephealthcheck"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	httpServer *httpserver.HttpServer
	tworker    worker.Worker
}

func NewWorker(ctx context.Context, cfg *WorkerConfig) (*Worker, error) {
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
		return nil, fmt.Errorf("error creating server %w", err)
	}

	httpServer.RegisterConnectHandler(func(connectServer *http.ServeMux) {
		connectServer.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	})

	// load secrets from gsm - not yet

	slog.InfoContext(ctx, "Initialising temporal client")
	tq := cfg.Temporal.TaskQueue
	tc, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Uri.AsString(),
		Namespace: cfg.Temporal.Namespace,
		Logger:    slog.Default(),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating temporal client %w", err)
	}

	slog.InfoContext(ctx, "Initialising temporal worker")
	w := worker.New(tc, tq, worker.Options{})

	// register workflows
	deephealthcheck.RegisterWorkflows(w)

	// register activities

	return &Worker{
		httpServer: httpServer,
		tworker:    w,
	}, nil
}

func (wk *Worker) Run(ctx context.Context) {
	// Run server
	rg := rungroup.NewRunGroup(ctx)
	rg.Run(func(ctx context.Context) {
		slog.InfoContext(ctx, "starting worker")
		if err := wk.tworker.Run(rungroup.InterruptChannel(ctx)); err != nil {
			slog.ErrorContext(ctx, "worker error", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "worker stopped")
	})
	rg.Run(func(ctx context.Context) {
		slog.InfoContext(ctx, "starting http server")
		if err := wk.httpServer.ListenAndServe(rungroup.InterruptChannel(ctx)); err != nil {
			slog.ErrorContext(ctx, "http server error", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "http server stopped")
	})
	rg.Wait()
}
