package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/grpchealth"
	"github.com/pseudo-su/golang-temporal-service-template/modules/frontdoor/internal/servicechecks"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1/deephealth_v1connect"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/httpserver"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/rungroup"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/pkg/workflowclient"
	"go.temporal.io/sdk/client"
)

type Frontdoor struct {
	httpServer *httpserver.HttpServer
}

func NewFrontdoor(ctx context.Context, cfg *FrontdoorConfig) (*Frontdoor, error) {
	slog.InfoContext(ctx, "Initialising temporal client")
	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Uri.AsString(),
		Namespace: cfg.Temporal.Namespace,
		Logger:    slog.Default(),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating temporal client %w", err)
	}

	slog.InfoContext(ctx, "Initialising workflow client")
	workflowClient, err := workflowclient.NewWorkflowClient(
		temporalClient,
		workflowclient.WithTaskQueue(cfg.Temporal.TaskQueue),
		workflowclient.WithNamespace(cfg.Temporal.Namespace),
	)
	if err != nil {
		return nil, fmt.Errorf("error workflow client %w", err)
	}

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

	deepHealthConnectServer := servicechecks.NewDeepHealthConnectServer(workflowClient)
	httpServer.RegisterConnectHandler(func(connectServer *http.ServeMux) {
		p, h := deephealth_v1connect.NewDeepHealthHandler(deepHealthConnectServer)
		connectServer.Handle(p, h)
		connectServer.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	})

	return &Frontdoor{
		httpServer: httpServer,
	}, nil
}

func (fd *Frontdoor) Run(ctx context.Context) {
	// Run server
	rg := rungroup.NewRunGroup(ctx)
	rg.Run(func(ctx context.Context) {
		slog.InfoContext(ctx, "starting http server")
		if err := fd.httpServer.ListenAndServe(rungroup.InterruptChannel(ctx)); err != nil {
			slog.ErrorContext(ctx, "http server error", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "http server stopped")
	})
	rg.Wait()
}
