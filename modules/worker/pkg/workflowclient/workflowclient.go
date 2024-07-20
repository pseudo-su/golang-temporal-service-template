package workflowclient

import (
	"context"
	"errors"
	"fmt"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	temporal "go.temporal.io/sdk/client"
)

var (
	ErrClientSetup    = errors.New("unable to setup workflow client")
	ErrTaskQueueEmpty = fmt.Errorf("task queue empty: %w", ErrClientSetup)
)

type clientOption func(c *WorkflowClient)

func WithTaskQueue(queue string) clientOption {
	return func(c *WorkflowClient) {
		c.taskQueue = queue
	}
}

func WithNamespace(namespace string) clientOption {
	return func(c *WorkflowClient) {
		c.namespace = namespace
	}
}

type WorkflowClientInterface interface {
	// Deep health check methods
	RunDeepHealthCheck(ctx context.Context, in *RunDeepHealthCheckInput) (*grpc_deephealth_v1.DeepHealthCheckResponse, error)
}

type WorkflowClient struct {
	temporalClient temporal.Client
	taskQueue      string
	namespace      string
}

func NewWorkflowClient(client temporal.Client, opts ...clientOption) (WorkflowClientInterface, error) {
	c := &WorkflowClient{
		temporalClient: client,
		taskQueue:      "",
		namespace:      "",
	}
	for _, optFn := range opts {
		optFn(c)
	}

	if c.taskQueue == "" {
		return nil, fmt.Errorf("%w", ErrTaskQueueEmpty)
	}

	return c, nil
}
