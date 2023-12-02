// Code generated by mockery v2.31.4. DO NOT EDIT.

package vendormocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	mock "github.com/stretchr/testify/mock"
)

// MockHealthClient is an autogenerated mock type for the HealthClient type
type MockHealthClient struct {
	mock.Mock
}

type MockHealthClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHealthClient) EXPECT() *MockHealthClient_Expecter {
	return &MockHealthClient_Expecter{mock: &_m.Mock}
}

// Check provides a mock function with given fields: ctx, in, opts
func (_m *MockHealthClient) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *grpc_health_v1.HealthCheckResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) *grpc_health_v1.HealthCheckResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc_health_v1.HealthCheckResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHealthClient_Check_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Check'
type MockHealthClient_Check_Call struct {
	*mock.Call
}

// Check is a helper method to define mock.On call
//   - ctx context.Context
//   - in *grpc_health_v1.HealthCheckRequest
//   - opts ...grpc.CallOption
func (_e *MockHealthClient_Expecter) Check(ctx interface{}, in interface{}, opts ...interface{}) *MockHealthClient_Check_Call {
	return &MockHealthClient_Check_Call{Call: _e.mock.On("Check",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockHealthClient_Check_Call) Run(run func(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption)) *MockHealthClient_Check_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*grpc_health_v1.HealthCheckRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockHealthClient_Check_Call) Return(_a0 *grpc_health_v1.HealthCheckResponse, _a1 error) *MockHealthClient_Check_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHealthClient_Check_Call) RunAndReturn(run func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error)) *MockHealthClient_Check_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, in, opts
func (_m *MockHealthClient) Watch(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 grpc_health_v1.Health_WatchClient
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) grpc_health_v1.Health_WatchClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc_health_v1.Health_WatchClient)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHealthClient_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type MockHealthClient_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - in *grpc_health_v1.HealthCheckRequest
//   - opts ...grpc.CallOption
func (_e *MockHealthClient_Expecter) Watch(ctx interface{}, in interface{}, opts ...interface{}) *MockHealthClient_Watch_Call {
	return &MockHealthClient_Watch_Call{Call: _e.mock.On("Watch",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockHealthClient_Watch_Call) Run(run func(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption)) *MockHealthClient_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*grpc_health_v1.HealthCheckRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockHealthClient_Watch_Call) Return(_a0 grpc_health_v1.Health_WatchClient, _a1 error) *MockHealthClient_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHealthClient_Watch_Call) RunAndReturn(run func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error)) *MockHealthClient_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHealthClient creates a new instance of MockHealthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHealthClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHealthClient {
	mock := &MockHealthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
