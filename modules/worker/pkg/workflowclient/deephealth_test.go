package workflowclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/testing-tools/vendormocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	errTestError = errors.New("something went wrong")
)

type TestWorkflowClient struct {
	suite.Suite
	ctx context.Context

	temporalClient *vendormocks.MockTemporalClient
	workflowClient WorkflowClientInterface
	workflowRunner *vendormocks.MockTemporalWorkflowRun
	result         *grpc_deephealth_v1.DeepHealthCheckResponse
}

func TestWorkflowClientSuite(t *testing.T) {
	suite.Run(t, &TestWorkflowClient{})
}

func (t *TestWorkflowClient) BeforeTest(_, _ string) {
	t.setupSuite()
}

func (t *TestWorkflowClient) SetupSubTest() {
	t.setupSuite()
}

func (t *TestWorkflowClient) setupSuite() {
	ctx := context.Background()
	c := vendormocks.NewMockTemporalClient(t.T())
	w, err := NewWorkflowClient(c, WithTaskQueue("test-task-queue"))
	require.NoError(t.T(), err)
	workflowRunner := vendormocks.NewMockTemporalWorkflowRun(t.T())

	t.workflowClient = w
	t.ctx = ctx
	t.result = &grpc_deephealth_v1.DeepHealthCheckResponse{}
	t.workflowRunner = workflowRunner
	t.temporalClient = c
}

func (t *TestWorkflowClient) TestWorkflowClient_RunDeepHealthCheck1() {
	tests := []struct {
		name    string
		mocks   func()
		want    *grpc_deephealth_v1.DeepHealthCheckResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "deep healthcheck workflow executes successfully",
			wantErr: assert.NoError,
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
				HttpStatus:  http.StatusOK,
			},
			mocks: func() {
				t.setupMocks(mockHappyWorkflowRun())
			},
		},
		{
			name:    "deep healthcheck workflow throws an error",
			wantErr: assert.Error,
			want:    nil,
			mocks: func() {
				t.setupMocks(mockErrorWorkflowRun())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func() {
			in := &RunDeepHealthCheckInput{}
			tt.mocks()
			got, err := t.workflowClient.RunDeepHealthCheck(t.ctx, in)
			if !tt.wantErr(t.T(), err, fmt.Sprintf("RunDeepHealthCheck(%v, %v)", t.ctx, in)) {
				return
			}
			assert.Equalf(t.T(), tt.want.String(), got.String(), "RunDeepHealthCheck(%v, %v)", t.ctx, in)
		})
	}
}

func (t *TestWorkflowClient) setupMocks(workflowResponse func(c context.Context, r interface{}) error) {
	t.temporalClient.EXPECT().ExecuteWorkflow(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(t.workflowRunner, nil)
	t.workflowRunner.
		EXPECT().
		Get(t.ctx, t.result).
		RunAndReturn(workflowResponse)
}

func mockHappyWorkflowRun() func(c context.Context, r interface{}) error {
	return func(c context.Context, r interface{}) error {
		res := r.(*grpc_deephealth_v1.DeepHealthCheckResponse)
		res.HealthState = grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK
		res.HttpStatus = http.StatusOK
		return nil
	}
}

func mockErrorWorkflowRun() func(c context.Context, r interface{}) error {
	return func(c context.Context, r interface{}) error {
		return errTestError
	}
}
