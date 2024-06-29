package httpserver

import (
	"context"
	"net/http"

	"connectrpc.com/grpcreflect"
)

type HttpServer struct {
	opts          httpServerOpts
	connectServer *http.ServeMux
}

func New(ctx context.Context, options ...Option) (*HttpServer, error) {
	s := &HttpServer{
		opts: httpServerOpts{},
	}
	for _, opt := range options {
		opt(s)
	}

	// TODO: how to setup interceptors
	// grpcOpts := []grpc.ServerOption{
	// 	grpc.ChainUnaryInterceptor(s.opts.unaryServerInterceptor...),
	// 	grpc.ChainStreamInterceptor(s.opts.streamServerInterceptor...),
	// }

	s.connectServer = http.NewServeMux()

	if s.opts.reflection {
		reflector := grpcreflect.NewStaticReflector()
		s.connectServer.Handle(grpcreflect.NewHandlerV1(reflector))
		s.connectServer.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}

	return s, nil
}

type RegisterConnectHandlerFn func(connectServer *http.ServeMux)

func (s *HttpServer) RegisterConnectHandler(registerFn RegisterConnectHandlerFn) {
	registerFn(s.connectServer)
}
