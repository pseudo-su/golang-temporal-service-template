package deephealthcheck

import (
	"context"
	"errors"
	"net/http"
	"testing"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/modules/worker/pkg/wftypes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

var (
	deepHealthCheckMethods = map[string]mockActivityMappingFixture{
		"MockDownstreamHealthCheckActivityOneFn": {fn: MockDownstreamHealthCheckActivityOneFn},
		"MockDownstreamHealthCheckActivityTwoFn": {fn: MockDownstreamHealthCheckActivityTwoFn},
	}
	errTestError = errors.New("something went wrong")
)

type DeepHealthTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func TestDeepHealthTestSuite(t *testing.T) {
	suite.Run(t, new(DeepHealthTestSuite))
}

func (s *DeepHealthTestSuite) BeforeTest(_, _ string) {
	s.env = s.NewTestWorkflowEnvironment()
	s.setAndRegisterHealthActivities()
}

func (s *DeepHealthTestSuite) SetupSubTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.setAndRegisterHealthActivities()
}

func (s *DeepHealthTestSuite) setAndRegisterHealthActivities() {
	s.T().Helper()
	Downstreams = []ServiceHealthActivityFn{
		MockDownstreamHealthCheckActivityOneFn,
		MockDownstreamHealthCheckActivityTwoFn,
	}
	s.env.RegisterActivity(MockDownstreamHealthCheckActivityOneFn)
	s.env.RegisterActivity(MockDownstreamHealthCheckActivityTwoFn)
}

func (s *DeepHealthTestSuite) AfterTest(_, _ string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeepHealthTestSuite) TearDownSubTest() {
	s.env.AssertExpectations(s.T())
}

func (s *DeepHealthTestSuite) Test_DeepHealthCheckWorkflow_SuccessNoErrors() {
	s.mockSuccessfulActivities(nil)
	s.env.ExecuteWorkflow(DeepHealthCheckWorkflow, nil)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result *grpc_deephealth_v1.DeepHealthCheckResponse
	err := s.env.GetWorkflowResult(&result)
	s.Nil(err)
	s.NotNil(result)

	s.Len(result.Services, 2, "expected two healthy services")
	for _, r := range result.Services {
		s.Equal(1, len(r.Checks))
		status, ok := r.Checks[0].Status.(*grpc_deephealth_v1.ServiceCheck_HttpStatus)
		s.True(ok)
		s.Equal(int32(http.StatusOK), status.HttpStatus)
	}
	s.Equal(int32(http.StatusOK), result.GetHttpStatus())
	s.Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK, result.GetHealthState())
}

func (s *DeepHealthTestSuite) Test_DeepHealthCheckWorkflow_Error_DueToDownstreamHealthFailure() {
	tests := []struct {
		name string
	}{
		{name: "MockDownstreamHealthCheckActivityOneFn"},
		{name: "MockDownstreamHealthCheckActivityTwoFn"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// mock all activities and ensure they will pass
			s.mockSuccessfulActivities(&tt.name)
			// make activity for test case fail
			s.mockDownstreamActivityErrorScenarios(tt.name, nil, false)

			s.env.ExecuteWorkflow(DeepHealthCheckWorkflow, nil)

			// workflow completes with no errors
			s.True(s.env.IsWorkflowCompleted())
			err := s.env.GetWorkflowError()
			s.NoError(err)

			// expect non nil result even with activity errors
			var result *grpc_deephealth_v1.DeepHealthCheckResponse
			err = s.env.GetWorkflowResult(&result)
			s.Nil(err)
			s.NotNil(result)

			s.Equal(int32(http.StatusServiceUnavailable), result.HttpStatus, "expected 503, got %v", result.HttpStatus)
			s.Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState, "expected %v, got %v", grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState)
			hasBadService := false
			for _, service := range result.Services {
				if len(service.Checks) != 1 {
					hasBadService = true
				}
				status, ok := service.Checks[0].Status.(*grpc_deephealth_v1.ServiceCheck_HttpStatus)
				if !ok || status.HttpStatus == http.StatusInternalServerError {
					hasBadService = true
				}
			}
			s.True(hasBadService, "at least one of the services is expected to have a non 200 status: %v", result.Services)
		})
	}
}

func (s *DeepHealthTestSuite) Test_DeepHealthCheckWorkflow_Error_DueToEmptyResponse() {
	tests := []struct {
		name string
	}{
		{name: "MockDownstreamHealthCheckActivityOneFn"},
		{name: "MockDownstreamHealthCheckActivityTwoFn"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// mock all activities and ensure they will pass
			s.mockSuccessfulActivities(&tt.name)
			// make activity for test case fail
			s.mockEmptyResponseActivities(&tt.name)

			s.env.ExecuteWorkflow(DeepHealthCheckWorkflow, nil)

			// workflow completes with no errors
			s.True(s.env.IsWorkflowCompleted())
			err := s.env.GetWorkflowError()
			s.NoError(err)

			// expect non nil result even with activity errors
			var result *grpc_deephealth_v1.DeepHealthCheckResponse
			err = s.env.GetWorkflowResult(&result)
			s.Nil(err)
			s.NotNil(result)

			s.Equal(int32(http.StatusServiceUnavailable), result.HttpStatus, "expected 503, got %v", result.HttpStatus)
			s.Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState, "expected %v, got %v", grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState)
		})
	}
}

func (s *DeepHealthTestSuite) Test_DeepHealthCheckWorkflow_Error_IgnoredDueToDisabled() {
	tests := []struct {
		name string
	}{
		{name: "MockDownstreamHealthCheckActivityOneFn"},
		{name: "MockDownstreamHealthCheckActivityTwoFn"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// mock all activities and ensure they will pass
			s.mockSuccessfulActivities(&tt.name)
			// make activity for test case fail
			s.mockDownstreamActivityErrorScenarios(tt.name, nil, true)

			s.env.ExecuteWorkflow(DeepHealthCheckWorkflow, nil)

			// workflow completes with no errors
			s.True(s.env.IsWorkflowCompleted())
			err := s.env.GetWorkflowError()
			s.NoError(err)

			// expect non nil result even with activity errors
			var result *grpc_deephealth_v1.DeepHealthCheckResponse
			err = s.env.GetWorkflowResult(&result)
			s.Nil(err)
			s.NotNil(result)

			s.Len(result.Services, 2, "expected two services")
			disabledService := 0
			healthyService := 0
			unhealthyService := 0
			for _, r := range result.Services {
				if r.ServiceMode == grpc_deephealth_v1.Service_SERVICE_MODE_DISABLED {
					disabledService++
				}

				if len(r.Checks) != 1 {
					unhealthyService++
					continue
				}

				if r.Checks[0].HealthState == grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_GOOD {
					healthyService++
				}

				if r.ServiceMode != grpc_deephealth_v1.Service_SERVICE_MODE_DISABLED && (r.Checks[0].HealthState == grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_BAD ||
					r.Checks[0].HealthState == grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_UNSPECIFIED) {
					unhealthyService++
				}
			}
			s.NotZero(disabledService)
			s.NotZero(healthyService)
			s.Zero(unhealthyService)
			s.Equal(int32(http.StatusOK), result.GetHttpStatus())
			s.Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK, result.GetHealthState())
		})
	}
}

func (s *DeepHealthTestSuite) Test_DeepHealthCheckWorkflow_Error_DueToDownstreamActivityError() {
	tests := []struct {
		name string
	}{
		{name: "MockDownstreamHealthCheckActivityOneFn"},
		{name: "MockDownstreamHealthCheckActivityTwoFn"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// mock all activities and ensure they will pass
			s.mockSuccessfulActivities(&tt.name)
			// make activity for test case fail
			s.mockDownstreamActivityErrorScenarios(tt.name, temporal.NewNonRetryableApplicationError("something went wrong", "activity-error", errTestError), false)

			s.env.ExecuteWorkflow(DeepHealthCheckWorkflow, nil)

			// workflow completes with no errors
			s.True(s.env.IsWorkflowCompleted())
			err := s.env.GetWorkflowError()
			s.NoError(err)

			// expect non nil result even with activity errors
			var result *grpc_deephealth_v1.DeepHealthCheckResponse
			err = s.env.GetWorkflowResult(&result)
			s.Nil(err)
			s.NotNil(result)

			s.Equal(int32(http.StatusServiceUnavailable), result.HttpStatus, "expected 503, got %v", result.HttpStatus)
			s.Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState, "expected %v, got %v", grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, result.HealthState)
		})
	}
}

func (s *DeepHealthTestSuite) mockSuccessfulActivities(exclusion *string) {
	s.T().Helper()
	s.mockDownstreamActivitySuccess(deepHealthCheckMethods, exclusion)
}

func (s *DeepHealthTestSuite) mockEmptyResponseActivities(exclusion *string) {
	s.T().Helper()
	s.mockDownstreamActivityEmptyResponse(deepHealthCheckMethods, exclusion)
}

type mockActivityMappingFixture struct {
	fn  func(ctx context.Context, in *wftypes.DeepHealthInput) (*grpc_deephealth_v1.Service, error)
	out *grpc_deephealth_v1.Service
	err error
}

func (s *DeepHealthTestSuite) mockDownstreamActivitySuccess(fixtures map[string]mockActivityMappingFixture, exclusion *string) {
	s.T().Helper()
	for key, fixture := range fixtures {
		if exclusion == nil || key != *exclusion {
			fixture.out = &grpc_deephealth_v1.Service{
				ServiceName: key,
				ServiceMode: grpc_deephealth_v1.Service_SERVICE_MODE_STUB,
				Checks: []*grpc_deephealth_v1.ServiceCheck{
					{
						HealthState: grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_GOOD,
						Status:      &grpc_deephealth_v1.ServiceCheck_HttpStatus{HttpStatus: http.StatusOK},
					},
				},
			}
			s.mockDownstreamActivity(fixture)
		}
	}
}
func (s *DeepHealthTestSuite) mockDownstreamActivityEmptyResponse(fixtures map[string]mockActivityMappingFixture, exclusion *string) {
	s.T().Helper()
	for key, fixture := range fixtures {
		if exclusion == nil || key == *exclusion {
			s.mockActivityEmptyResponse(fixture)
		}
	}
}

func (s *DeepHealthTestSuite) mockDownstreamActivityErrorScenarios(name string, err error, isDisabled bool) {
	s.T().Helper()
	out := &grpc_deephealth_v1.Service{
		ServiceName: name,
		ServiceMode: grpc_deephealth_v1.Service_SERVICE_MODE_STUB,
		Checks: []*grpc_deephealth_v1.ServiceCheck{
			{
				HealthState: grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_BAD,
				Status:      &grpc_deephealth_v1.ServiceCheck_HttpStatus{HttpStatus: http.StatusInternalServerError},
			},
		},
	}
	if isDisabled {
		out.ServiceMode = grpc_deephealth_v1.Service_SERVICE_MODE_DISABLED
	}
	s.mockDownstreamActivity(mockActivityMappingFixture{
		fn:  deepHealthCheckMethods[name].fn,
		out: out,
		err: err,
	})
}

func (s *DeepHealthTestSuite) mockDownstreamActivity(fixture mockActivityMappingFixture) {
	s.T().Helper()
	s.env.OnActivity(fixture.fn, mock.Anything, mock.Anything).
		Return(fixture.out, fixture.err).
		Times(1)
}

func (s *DeepHealthTestSuite) mockActivityEmptyResponse(fixture mockActivityMappingFixture) {
	s.T().Helper()
	s.env.OnActivity(fixture.fn, mock.Anything, mock.Anything).
		Return(nil, nil).
		Times(1)
}

func MockDownstreamHealthCheckActivityOneFn(ctx context.Context, in *wftypes.DeepHealthInput) (*grpc_deephealth_v1.Service, error) {
	// this can be replaced by actual activities at some point, but given it's all just mocking it makes little difference
	return &grpc_deephealth_v1.Service{}, nil
}

func MockDownstreamHealthCheckActivityTwoFn(ctx context.Context, in *wftypes.DeepHealthInput) (*grpc_deephealth_v1.Service, error) {
	// this can be replaced by actual activities at some point, but given it's all just mocking it makes little difference
	return &grpc_deephealth_v1.Service{}, nil
}
