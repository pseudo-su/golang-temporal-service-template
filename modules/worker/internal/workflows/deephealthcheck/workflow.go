package deephealthcheck

import (
	"context"
	"time"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/servicechecks"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/internal/workflowutils"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/pkg/wftypes"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceHealthActivityFn func(ctx context.Context, in *wftypes.DeepHealthInput) (*grpc_deephealth_v1.Service, error)

// Downstreams corresponding activities while able to return an error should avoid returning errors and rather utilise
// corresponding values from grpc_deephealth_v1.Service to return a failed result, excepting the instance where
// when service ENABLED=false, corresponding health check should still happen but due to being disabled should not make the overall deep health check fail
var Downstreams = []ServiceHealthActivityFn{}

func RegisterWorkflows(w worker.WorkflowRegistry) {
	w.RegisterWorkflow(DeepHealthCheckWorkflow)
}

func DeepHealthCheckWorkflow(ctx workflow.Context, in *wftypes.DeepHealthInput) (*grpc_deephealth_v1.DeepHealthCheckResponse, error) {
	// setup default values for response
	ctx = workflow.WithStartToCloseTimeout(ctx, time.Second*10)
	logger := workflow.GetLogger(ctx)

	logger.Debug("starting DeepHealthCheckWorkflow")

	var activityInput *wftypes.DeepHealthInput

	if in != nil && in.User != nil {
		activityInput = &wftypes.DeepHealthInput{
			User: in.User,
		}
	}

	selector := workflow.NewSelector(ctx)
	aoCtx := workflowutils.WithNoRetriesRetryPolicy(ctx, nil)
	var res []servicechecks.HealthActivityResponse
	// start the activities, add them to the selector as a future
	for _, downstreamFn := range Downstreams {
		activityFn := downstreamFn
		workflowFuture := workflow.ExecuteActivity(aoCtx, activityFn, activityInput)
		selector.AddFuture(workflowFuture, func(f workflow.Future) {
			var result *grpc_deephealth_v1.Service
			err := f.Get(ctx, &result)
			res = append(res, servicechecks.HealthActivityResponse{
				Result: result,
				Err:    err,
			})
		})
	}

	// run through selector and as future are made callable retrieve from it
	for range Downstreams {
		selector.Select(ctx)
	}

	serviceCheckResponse, err := servicechecks.ProcessDeepHealthChecks(res, grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER, len(Downstreams))
	if err != nil {
		logger.Error("error on deep health check", "error", err, "response", serviceCheckResponse)
	}
	serviceCheckResponse.GeneratedTime = timestamppb.New(workflow.Now(ctx).UTC())
	return serviceCheckResponse, nil
}
