package servicechecks

import (
	"context"

	"connectrpc.com/connect"
	deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1/deephealth_v1connect"
)

type DeepHealthConnectServer struct {
	deephealth_v1connect.UnimplementedDeepHealthHandler
}

var _ deephealth_v1connect.DeepHealthHandler = &DeepHealthConnectServer{}

func NewDeepHealthConnectServer() *DeepHealthConnectServer {
	return &DeepHealthConnectServer{}
}

func (DeepHealthConnectServer) Check(ctx context.Context, req *connect.Request[deephealth_v1.CheckRequest]) (*connect.Response[deephealth_v1.CheckResponse], error) {
	return connect.NewResponse(
		&deephealth_v1.CheckResponse{
			Message: "Placeholder",
		},
	), nil
}
