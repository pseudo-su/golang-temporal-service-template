package workflowclient

import (
	"context"
	"time"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/internal/workflows/deephealthcheck"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/pkg/wftypes"
	s_temporal "go.temporal.io/sdk/temporal"

	temporal "go.temporal.io/sdk/client"
)

type RunDeepHealthCheckInput struct {
	TriggeredBy string
	User        *wftypes.DeepHealthUser
}

func (c *WorkflowClient) RunDeepHealthCheck(ctx context.Context, in *RunDeepHealthCheckInput) (*grpc_deephealth_v1.DeepHealthCheckResponse, error) {
	searchAttribute := s_temporal.NewSearchAttributeKeyKeyword("HealthCheckTriggeredBy").ValueSet(in.TriggeredBy)
	opts := temporal.StartWorkflowOptions{
		TaskQueue:                c.taskQueue,
		TypedSearchAttributes:    s_temporal.NewSearchAttributes(searchAttribute),
		WorkflowExecutionTimeout: time.Second * 60,
	}
	resultFuture, err := c.temporalClient.ExecuteWorkflow(ctx, opts, deephealthcheck.DeepHealthCheckWorkflow, &wftypes.DeepHealthInput{
		User: in.User,
	})
	if err != nil {
		return nil, err
	}

	result := &grpc_deephealth_v1.DeepHealthCheckResponse{}
	err = resultFuture.Get(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
