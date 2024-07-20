package servicechecks

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/stretchr/testify/assert"
)

var (
	errTest                 = errors.New("something went wrong")
	healthyActivityResponse = []HealthActivityResponse{{
		Result: mockGoodService("service-1", "check-1"),
		Err:    nil,
	}}

	unhealthyActivityResponse = []HealthActivityResponse{{
		Result: mockBadService("service-1", "check-1"),
		Err:    nil,
	}}

	multipleResponseWithOneError = []HealthActivityResponse{
		{
			Result: mockGoodService("service-1", "check-1"),
			Err:    nil,
		},
		{
			Result: nil,
			Err:    errTest,
		},
	}
)

func TestProcessDeepHealthChecks(t *testing.T) {
	type args struct {
		responses              []HealthActivityResponse
		initFrom               grpc_deephealth_v1.ServiceCheck_InitiatedFromService
		expectedTotalResponses int
	}
	tests := []struct {
		name    string
		args    args
		want    *grpc_deephealth_v1.DeepHealthCheckResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty health activity response with 0 expected total responses",
			args: args{
				responses:              []HealthActivityResponse{},
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 0,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
				HttpStatus:  http.StatusNotImplemented,
				Services:    nil,
			},
			wantErr: assertError(errNoServicesConfigured),
		},
		{
			name: "health activity response containing one error and one valid response and expecting 2 total responses - contains case specific error",
			args: args{
				responses:              multipleResponseWithOneError,
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 2,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL,
				HttpStatus:  http.StatusServiceUnavailable,
				Services:    []*grpc_deephealth_v1.Service{mockGoodService("service-1", "check-1")},
			},
			wantErr: assertError(errTotalResponse),
		},
		{
			name: "health activity response containing one error and one valid response and expecting 2 total responses - contains original error",
			args: args{
				responses:              multipleResponseWithOneError,
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 2,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL,
				HttpStatus:  http.StatusServiceUnavailable,
				Services:    []*grpc_deephealth_v1.Service{mockGoodService("service-1", "check-1")},
			},
			wantErr: assertError(errTest),
		},
		{
			name: "health activity response good with 1 expected total responses",
			args: args{
				responses:              healthyActivityResponse,
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 1,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
				HttpStatus:  http.StatusOK,
				Services:    []*grpc_deephealth_v1.Service{mockGoodService("service-1", "check-1")},
			},
			wantErr: assert.NoError,
		},
		{
			name: "health activity response 1 good 1 bad but disabled with 2 expected total responses",
			args: args{
				responses: []HealthActivityResponse{
					{Result: mockGoodService("service-1", "check-1"), Err: nil},
					{Result: mockDisabledBadService("service-2", "check-1"), Err: nil},
				},
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 2,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
				HttpStatus:  http.StatusOK,
				Services: []*grpc_deephealth_v1.Service{
					mockGoodService("service-1", "check-1"),
					mockDisabledBadService("service-2", "check-1"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "health activity 2 checks for the same service get merged",
			args: args{
				responses: []HealthActivityResponse{
					{Result: mockGoodService("service-1", "check-1"), Err: nil},
					{Result: mockGoodService("service-1", "check-2"), Err: nil},
				},
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 2,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
				HttpStatus:  http.StatusOK,
				Services: []*grpc_deephealth_v1.Service{
					{
						ServiceName: "service-1",
						ServiceMode: grpc_deephealth_v1.Service_SERVICE_MODE_LIVE,
						Checks: []*grpc_deephealth_v1.ServiceCheck{
							{
								Description:   "check-1",
								Status:        &grpc_deephealth_v1.ServiceCheck_HttpStatus{HttpStatus: http.StatusOK},
								HealthState:   grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_GOOD,
								InitiatedFrom: grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
							},
							{
								Description:   "check-2",
								Status:        &grpc_deephealth_v1.ServiceCheck_HttpStatus{HttpStatus: http.StatusOK},
								HealthState:   grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_GOOD,
								InitiatedFrom: grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "health activity response bad with 1 expected total response",
			args: args{
				responses:              unhealthyActivityResponse,
				initFrom:               grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
				expectedTotalResponses: 1,
			},
			want: &grpc_deephealth_v1.DeepHealthCheckResponse{
				HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL,
				HttpStatus:  http.StatusServiceUnavailable,
				Services:    []*grpc_deephealth_v1.Service{mockBadService("service-1", "check-1")},
			},
			wantErr: assertError(errServiceInFailedState),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessDeepHealthChecks(tt.args.responses, tt.args.initFrom, tt.args.expectedTotalResponses)
			tt.wantErr(t, err, fmt.Sprintf("ProcessDeepHealthChecks(%v, %v, %v)", tt.args.responses, tt.args.initFrom, tt.args.expectedTotalResponses))
			assert.Equalf(t, tt.want.String(), got.String(), "ProcessDeepHealthChecks(%v, %v, %v)", tt.args.responses, tt.args.initFrom, tt.args.expectedTotalResponses)
		})
	}
}

func assertError(expectedErr error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		if err == nil || !errors.Is(err, expectedErr) {
			return assert.Fail(t, fmt.Sprintf("expected %v - got %v", expectedErr, err))
		}
		return true
	}
}

func mockGoodService(serviceName, description string) *grpc_deephealth_v1.Service {
	return mockService(serviceName, description, grpc_deephealth_v1.Service_SERVICE_MODE_LIVE, grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_GOOD, http.StatusOK)
}

func mockDisabledBadService(serviceName, description string) *grpc_deephealth_v1.Service {
	return mockService(serviceName, description, grpc_deephealth_v1.Service_SERVICE_MODE_DISABLED, grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_BAD, http.StatusInternalServerError)
}

func mockBadService(serviceName, description string) *grpc_deephealth_v1.Service {
	return mockService(serviceName, description, grpc_deephealth_v1.Service_SERVICE_MODE_LIVE, grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_BAD, http.StatusInternalServerError)
}

func mockService(serviceName, description string, serviceMode grpc_deephealth_v1.Service_ServiceMode, healthState grpc_deephealth_v1.ServiceCheck_HealthState, httpStatus int32) *grpc_deephealth_v1.Service {
	return &grpc_deephealth_v1.Service{
		ServiceName: serviceName,
		ServiceMode: serviceMode,
		Checks: []*grpc_deephealth_v1.ServiceCheck{
			{
				Description:   description,
				Status:        &grpc_deephealth_v1.ServiceCheck_HttpStatus{HttpStatus: http.StatusOK},
				HealthState:   healthState,
				InitiatedFrom: grpc_deephealth_v1.ServiceCheck_INITIATED_FROM_SERVICE_WORKER,
			},
		},
	}
}
