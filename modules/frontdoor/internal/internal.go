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
)

type Frontdoor struct {
	httpServer *httpserver.HttpServer
}

func NewFrontdoor(ctx context.Context, cfg *FrontdoorConfig) (*Frontdoor, error) {
	slog.InfoContext(ctx, "Initialising http server")
	address := fmt.Sprintf(":%d", cfg.Tcp.Port)
	httpServer, err := httpserver.New(
		ctx,
		httpserver.WithAddress(address),
		httpserver.WithReflection(),
		httpserver.WithUnaryInterceptors(),
		httpserver.WithStreamInterceptors(),
		// httpserver.WithEmbeddedGrpcGateway(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating server %w", err)
	}

	// healthServer := health.NewServer()
	// deepHealthGrpcServer := servicechecks.NewDeepHealthGrpcServer()
	deepHealthConnectServer := servicechecks.NewDeepHealthConnectServer()
	httpServer.RegisterConnectHandler(func(connectServer *http.ServeMux) {
		p, h := deephealth_v1connect.NewDeepHealthHandler(deepHealthConnectServer)
		connectServer.Handle(p, h)
		connectServer.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))
	})
	// httpServer.RegisterGrpcServer(func(s *grpc.Server) {
	// 	grpc_health_v1.RegisterHealthServer(s, healthServer)
	// 	grpc_deephealth_v1.RegisterDeepHealthServer(s, deepHealthGrpcServer)
	// })
	// err = httpServer.RegisterGatewayHandlers(func(grpcGatewayServer *runtime.ServeMux) error {
	// 	opts := []grpc.DialOption{
	// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	}
	// 	return errors.Join(
	// 		grpc_deephealth_v1.RegisterDeepHealthHandlerFromEndpoint(ctx, grpcGatewayServer, address, opts),
	// 	)
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("error registering gateway handlers %w", err)
	// }
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
