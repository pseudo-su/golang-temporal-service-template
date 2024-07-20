package servicechecks

import (
	"context"

	"connectrpc.com/connect"
	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1/deephealth_v1connect"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/pkg/workflowclient"
)

type DeepHealthCheckerInterface interface {
	RunDeepHealthCheck(ctx context.Context, in *workflowclient.RunDeepHealthCheckInput) (*grpc_deephealth_v1.DeepHealthCheckResponse, error)
}

type DeepHealthConnectServer struct {
	deephealth_v1connect.UnimplementedDeepHealthHandler

	workflowChecker DeepHealthCheckerInterface
	// frontdoorChecker DeepHealthCheckerInterface
}

var _ deephealth_v1connect.DeepHealthHandler = &DeepHealthConnectServer{}

func NewDeepHealthConnectServer(workflowClient workflowclient.WorkflowClientInterface) *DeepHealthConnectServer {
	return &DeepHealthConnectServer{
		workflowChecker: workflowClient,
	}
}

func (s *DeepHealthConnectServer) Check(ctx context.Context, req *connect.Request[grpc_deephealth_v1.DeepHealthCheckRequest]) (*connect.Response[grpc_deephealth_v1.DeepHealthCheckResponse], error) {
	workflowServices, err := s.workflowChecker.RunDeepHealthCheck(ctx, &workflowclient.RunDeepHealthCheckInput{
		TriggeredBy: req.Msg.GetTriggeredBy(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(workflowServices), nil
}
