// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: deephealth/v1/deephealth.proto

/*
Package deephealth_v1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package deephealth_v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_DeepHealth_Check_0(ctx context.Context, marshaler runtime.Marshaler, client DeepHealthClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CheckRequest
	var metadata runtime.ServerMetadata

	msg, err := client.Check(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_DeepHealth_Check_0(ctx context.Context, marshaler runtime.Marshaler, server DeepHealthServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CheckRequest
	var metadata runtime.ServerMetadata

	msg, err := server.Check(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterDeepHealthHandlerServer registers the http handlers for service DeepHealth to "mux".
// UnaryRPC     :call DeepHealthServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterDeepHealthHandlerFromEndpoint instead.
func RegisterDeepHealthHandlerServer(ctx context.Context, mux *runtime.ServeMux, server DeepHealthServer) error {

	mux.Handle("GET", pattern_DeepHealth_Check_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/deephealth.v1.DeepHealth/Check", runtime.WithHTTPPathPattern("/v1/health/deep"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_DeepHealth_Check_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_DeepHealth_Check_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterDeepHealthHandlerFromEndpoint is same as RegisterDeepHealthHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterDeepHealthHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterDeepHealthHandler(ctx, mux, conn)
}

// RegisterDeepHealthHandler registers the http handlers for service DeepHealth to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterDeepHealthHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterDeepHealthHandlerClient(ctx, mux, NewDeepHealthClient(conn))
}

// RegisterDeepHealthHandlerClient registers the http handlers for service DeepHealth
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "DeepHealthClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "DeepHealthClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "DeepHealthClient" to call the correct interceptors.
func RegisterDeepHealthHandlerClient(ctx context.Context, mux *runtime.ServeMux, client DeepHealthClient) error {

	mux.Handle("GET", pattern_DeepHealth_Check_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/deephealth.v1.DeepHealth/Check", runtime.WithHTTPPathPattern("/v1/health/deep"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_DeepHealth_Check_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_DeepHealth_Check_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_DeepHealth_Check_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "health", "deep"}, ""))
)

var (
	forward_DeepHealth_Check_0 = runtime.ForwardResponseMessage
)
