package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/grpcreflect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type HttpServer struct {
	opts              httpServerOpts
	grpcServer        *grpc.Server
	grpcGatewayServer *runtime.ServeMux
	connectServer     *http.ServeMux
	combinedServer    *http.Server
}

func New(ctx context.Context, options ...Option) (*HttpServer, error) {
	s := &HttpServer{
		opts: httpServerOpts{},
	}
	for _, opt := range options {
		opt(s)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(s.opts.unaryServerInterceptor...),
		grpc.ChainStreamInterceptor(s.opts.streamServerInterceptor...),
	}

	s.grpcServer = grpc.NewServer(grpcOpts...)

	if s.opts.reflection {
		reflection.Register(s.grpcServer)
	}

	if s.opts.useEmbeddedGrpcGateway {
		s.grpcGatewayServer = runtime.NewServeMux()
	}

	s.connectServer = http.NewServeMux()

	if s.opts.reflection {
		reflector := grpcreflect.NewStaticReflector()
		s.connectServer.Handle(grpcreflect.NewHandlerV1(reflector))
		s.connectServer.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}

	s.combinedServer = newCombinedServer(ctx, s.grpcServer, s.grpcGatewayServer)

	return s, nil
}

type RegisterGrpcServerFn func(*grpc.Server)

func (s *HttpServer) RegisterGrpcServer(registerFn RegisterGrpcServerFn) {
	registerFn(s.grpcServer)
}

type RegisterGrpcGatewayHandlerFn func(grpcGatewayServer *runtime.ServeMux) error

func (s *HttpServer) RegisterGatewayHandlers(registerFn RegisterGrpcGatewayHandlerFn) error {
	if !s.opts.useEmbeddedGrpcGateway {
		return fmt.Errorf("attempted to register gateway handlers without enabling embedded gateway")
	}
	return registerFn(s.grpcGatewayServer)
}

type RegisterConnectHandlerFn func(connectServer *http.ServeMux)

func (s *HttpServer) RegisterConnectHandler(registerFn RegisterConnectHandlerFn) {
	registerFn(s.connectServer)
}
