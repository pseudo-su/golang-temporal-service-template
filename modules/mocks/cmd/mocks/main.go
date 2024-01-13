package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pseudo-su/golang-temporal-service-template/modules/mocks/internal"
	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/envconfig"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/httpserver"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/logsetup"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/rungroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var cfg *internal.MocksConfig

func init() {
	ctx := context.Background()

	cfg = envconfig.ParseEnv(&internal.MocksConfig{})
	slog.InfoContext(ctx, "App config loaded", slog.Any("name", cfg.App.Name), slog.Any("env", cfg.App.Env))

	logger := logsetup.NewLogger(
		logsetup.WithModeStr(cfg.Log.Mode),
		logsetup.WithLevelStr(cfg.Log.Level),
	)
	slog.SetDefault(logger)
}

func main() {
	ctx := context.Background()

	slog.InfoContext(ctx, "Initialising http server")
	address := fmt.Sprintf(":%d", cfg.Tcp.Port)
	httpServer, err := httpserver.New(
		ctx,
		httpserver.WithAddress(address),
		httpserver.WithReflection(),
		httpserver.WithUnaryInterceptors(),
		httpserver.WithStreamInterceptors(),
		httpserver.WithEmbeddedGrpcGateway(),
	)
	if err != nil {
		slog.ErrorContext(ctx, "error creating server", slog.Any("error", err))
		os.Exit(1)
	}

	healthServer := health.NewServer()
	httpServer.RegisterGrpcServer(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, healthServer)
	})
	err = httpServer.RegisterGatewayHandlers(func(grpcGatewayServer *runtime.ServeMux) error {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		return errors.Join(
			grpc_deephealth_v1.RegisterDeepHealthHandlerFromEndpoint(ctx, grpcGatewayServer, address, opts),
		)
	})
	if err != nil {
		slog.ErrorContext(ctx, "error registering gateway handlers", slog.Any("error", err))
		os.Exit(1)
	}

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
