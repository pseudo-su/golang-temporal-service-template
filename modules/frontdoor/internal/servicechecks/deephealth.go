package servicechecks

import (
	"context"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
)

type DeepHealthServer struct {
	grpc_deephealth_v1.UnimplementedDeepHealthServer
}

func NewDeepHealthServer() *DeepHealthServer {
	return &DeepHealthServer{}
}

var _ grpc_deephealth_v1.DeepHealthServer = &DeepHealthServer{}

func (s *DeepHealthServer) Check(context.Context, *grpc_deephealth_v1.CheckRequest) (*grpc_deephealth_v1.CheckResponse, error) {
	return &grpc_deephealth_v1.CheckResponse{
		Message: "Placeholder",
	}, nil
}
