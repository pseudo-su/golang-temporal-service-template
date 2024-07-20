package servicechecks

import (
	"errors"
	"fmt"
	"net/http"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
)

var (
	errEmptyServiceResponse = errors.New("received empty response from downstream")
	errNoServicesConfigured = errors.New("no services configured to report health on")
	errTotalResponse        = errors.New("number of service responses doesn't match total responses expected")
	errServiceInFailedState = errors.New("service health output in invalid state")
)

type HealthActivityResponse struct {
	Result *grpc_deephealth_v1.Service
	Err    error
}

// ProcessDeepHealthChecks intentionally returns no error as grpc_deephealth_v1.DeepHealthCheckResponse should reflect whether deep health check failed or not
func ProcessDeepHealthChecks(responses []HealthActivityResponse, initFrom grpc_deephealth_v1.ServiceCheck_InitiatedFromService, expectedTotalResponses int) (*grpc_deephealth_v1.DeepHealthCheckResponse, error) {
	var errList error
	serviceCheckResponse := &grpc_deephealth_v1.DeepHealthCheckResponse{
		HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK,
		HttpStatus:  http.StatusOK,
	}
	for _, response := range responses {
		errList = errors.Join(errList, recordErrorsAndAppendServiceResponse(serviceCheckResponse, response.Result, response.Err, initFrom))
	}
	errList = errors.Join(errList, updateCheckResponse(serviceCheckResponse, expectedTotalResponses, errList))
	return serviceCheckResponse, errList
}

func recordErrorsAndAppendServiceResponse(existing *grpc_deephealth_v1.DeepHealthCheckResponse, resultToProcess *grpc_deephealth_v1.Service, errToProcess error, initiatedFrom grpc_deephealth_v1.ServiceCheck_InitiatedFromService) error {
	if resultToProcess == nil {
		return errors.Join(errToProcess, errEmptyServiceResponse)
	}

	if resultToProcess.ServiceMode != grpc_deephealth_v1.Service_SERVICE_MODE_DISABLED {
		for _, check := range resultToProcess.Checks {
			if check.HealthState == grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_BAD ||
				check.HealthState == grpc_deephealth_v1.ServiceCheck_HEALTH_STATE_UNSPECIFIED {
				errToProcess = errors.Join(errToProcess, fmt.Errorf("service %s is in invalid state %s - %w", resultToProcess.ServiceName, check.HealthState, errServiceInFailedState))
				break
			}
		}
	}

	for _, check := range resultToProcess.Checks {
		check.InitiatedFrom = initiatedFrom
	}

	// append the Result to the list of downstream results
	Merge(existing, resultToProcess)

	return errToProcess
}

// If multiple checks are made against the same service, this function merges the checks under the same top level service.
// Example (pseudo-code):
//
// existing:
//
//	{
//	    Services: [
//	        {
//	            ServiceName: "forgerock",
//	            Checks: [ { Description: "jwks" } ]
//	        },
//	    ]
//	}
//
// service:
//
//	{
//	    ServiceName: "forgerock",
//	    Checks: [ { Description: "customer rqf token" } ]
//	},
//
// result:
//
//	{
//	    Services: [
//	        {
//	            ServiceName: "forgerock",
//	            Checks: [ { Description: "jwks" }, { Description: "customer rqf token" } ]
//	        },
//	    ]
//	}
func Merge(existing *grpc_deephealth_v1.DeepHealthCheckResponse, service *grpc_deephealth_v1.Service) {
	foundExisting := false
	for _, existingService := range existing.Services {
		if existingService.ServiceName != service.ServiceName {
			continue
		}

		existingService.Checks = append(existingService.Checks, service.Checks...)
		foundExisting = true
		break
	}
	if !foundExisting {
		existing.Services = append(existing.Services, service)
	}
}

func updateCheckResponse(serviceCheckResponse *grpc_deephealth_v1.DeepHealthCheckResponse, numExpectedResponses int, processedErrors error) error {
	// This sets the response state to failed and status unimplemented if no downstream calls were made
	var errResponse error
	if numExpectedResponses == 0 {
		serviceCheckResponse.HealthState = grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK
		serviceCheckResponse.HttpStatus = http.StatusNotImplemented
		errResponse = errNoServicesConfigured
	}
	// conditions to consider where overall DeepHealthCheckResponse has failed
	numChecks := 0
	for _, service := range serviceCheckResponse.Services {
		numChecks += len(service.Checks)
	}

	failedServiceCheckResponse := numChecks != numExpectedResponses || processedErrors != nil
	if failedServiceCheckResponse {
		serviceCheckResponse.HealthState = grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL
		serviceCheckResponse.HttpStatus = http.StatusServiceUnavailable
		errResponse = errTotalResponse
	}

	return errResponse
}
