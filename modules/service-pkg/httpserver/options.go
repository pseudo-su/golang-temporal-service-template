package httpserver

import (
	"net"

	"google.golang.org/grpc"
)

type Option func(*HttpServer)

type httpServerOpts struct {
	address                 string
	listener                net.Listener
	unaryServerInterceptor  []grpc.UnaryServerInterceptor
	streamServerInterceptor []grpc.StreamServerInterceptor
	reflection              bool
	useEmbeddedGrpcGateway  bool
}

func WithAddress(address string) Option {
	return func(s *HttpServer) {
		s.opts.address = address
	}
}

func WithListener(listener net.Listener) Option {
	return func(s *HttpServer) {
		s.opts.listener = listener
	}
}

func WithUnaryInterceptors(inter ...grpc.UnaryServerInterceptor) Option {
	return func(b *HttpServer) {
		b.opts.unaryServerInterceptor = inter
	}
}

func WithStreamInterceptors(inter ...grpc.StreamServerInterceptor) Option {
	return func(b *HttpServer) {
		b.opts.streamServerInterceptor = inter
	}
}

func WithReflection() Option {
	return func(b *HttpServer) {
		b.opts.reflection = true
	}
}

func WithEmbeddedGrpcGateway() Option {
	return func(s *HttpServer) {
		s.opts.useEmbeddedGrpcGateway = true
	}
}
