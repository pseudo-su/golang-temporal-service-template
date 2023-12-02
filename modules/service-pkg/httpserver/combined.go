package httpserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func isGrpcRequest(r *http.Request) bool {
	contentTypeHeader := r.Header.Get("Content-Type")
	return r.ProtoMajor == 2 && strings.HasPrefix(contentTypeHeader, "application/grpc")
}

func newCombinedServer(ctx context.Context, grpcServer *grpc.Server, grpcGatewayServer *runtime.ServeMux) *http.Server {
	return &http.Server{
		Handler: combineHandlers(
			ctx,
			grpcServer,
			grpcGatewayServer,
		),
	}
}

func combineHandlers(ctx context.Context, grpcServer *grpc.Server, grpcGatewayServer *runtime.ServeMux) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)

		if grpcGatewayServer != nil && !isGrpcRequest(req) {
			grpcGatewayServer.ServeHTTP(w, req)
		} else {
			grpcServer.ServeHTTP(w, req)
		}
	}

	return h2c.NewHandler(http.HandlerFunc(hf), &http2.Server{})
}
