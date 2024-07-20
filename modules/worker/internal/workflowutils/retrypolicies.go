package workflowutils

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RetryPolicies struct {
	InitialInterval    time.Duration
	BackOffCoefficient float64
	MaximumInterval    time.Duration
	MaximumAttempts    int32
	MaxRetryPeriod     time.Duration
}

var (
	c1        = RetryPolicies{time.Second, 1.2, 300 * time.Second, 122, 8 * time.Hour}
	c2        = RetryPolicies{time.Second, 1.2, 900 * time.Second, 224, 48 * time.Hour}
	c3        = RetryPolicies{time.Second, 1.2, 2700 * time.Second, 230, 144 * time.Hour}
	c4        = RetryPolicies{time.Second, 1.2, 5400 * time.Second, 362, 481 * time.Hour}
	c5        = RetryPolicies{time.Second, 1.2, 10800 * time.Second, 379, 1002 * time.Hour}
	noRetries = RetryPolicies{time.Second, 1.2, time.Second * 60, 1, time.Second * 60}
)

// WithC1RetryPolicy returns a context with an activity retry policy suitable for C1 services.
func WithC1RetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return c1.RetryPolicy(ctx, responseTime)
}

// WithC2RetryPolicy returns a context with an activity retry policy suitable for C2 services.
func WithC2RetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return c2.RetryPolicy(ctx, responseTime)
}

// WithC3RetryPolicy returns a context with an activity retry policy suitable for C3 services.
func WithC3RetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return c3.RetryPolicy(ctx, responseTime)
}

// WithC4RetryPolicy returns a context with an activity retry policy suitable for C4 services.
func WithC4RetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return c4.RetryPolicy(ctx, responseTime)
}

// WithC5RetryPolicy returns a context with an activity retry policy suitable for C5 services.
func WithC5RetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return c5.RetryPolicy(ctx, responseTime)
}

func WithNoRetriesRetryPolicy(ctx workflow.Context, responseTime *time.Duration) workflow.Context {
	return noRetries.RetryPolicy(ctx, responseTime)
}

// RetryPolicy returns a context with an activity retry policy suitable for C1, C2, C3, C4, and C5 services.
func (p RetryPolicies) RetryPolicy(ctx workflow.Context, averageExpectedResponseTime *time.Duration) workflow.Context {
	startToCloseTimeout := time.Minute * 2
	if averageExpectedResponseTime != nil {
		startToCloseTimeout = *averageExpectedResponseTime
	}
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: p.MaxRetryPeriod,
		StartToCloseTimeout:    startToCloseTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient:     p.BackOffCoefficient,
			InitialInterval:        p.InitialInterval,
			MaximumInterval:        p.MaximumInterval,
			MaximumAttempts:        p.MaximumAttempts,
			NonRetryableErrorTypes: []string{"non-retryable"},
		},
	})
}
