// Code generated by mockery v2.31.4. DO NOT EDIT.

package vendormocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	sdkclient "go.temporal.io/sdk/client"
)

// MockTemporalWorkflowRun is an autogenerated mock type for the WorkflowRun type
type MockTemporalWorkflowRun struct {
	mock.Mock
}

type MockTemporalWorkflowRun_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTemporalWorkflowRun) EXPECT() *MockTemporalWorkflowRun_Expecter {
	return &MockTemporalWorkflowRun_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, valuePtr
func (_m *MockTemporalWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	ret := _m.Called(ctx, valuePtr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, valuePtr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemporalWorkflowRun_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockTemporalWorkflowRun_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - valuePtr interface{}
func (_e *MockTemporalWorkflowRun_Expecter) Get(ctx interface{}, valuePtr interface{}) *MockTemporalWorkflowRun_Get_Call {
	return &MockTemporalWorkflowRun_Get_Call{Call: _e.mock.On("Get", ctx, valuePtr)}
}

func (_c *MockTemporalWorkflowRun_Get_Call) Run(run func(ctx context.Context, valuePtr interface{})) *MockTemporalWorkflowRun_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *MockTemporalWorkflowRun_Get_Call) Return(_a0 error) *MockTemporalWorkflowRun_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemporalWorkflowRun_Get_Call) RunAndReturn(run func(context.Context, interface{}) error) *MockTemporalWorkflowRun_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetID provides a mock function with given fields:
func (_m *MockTemporalWorkflowRun) GetID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockTemporalWorkflowRun_GetID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetID'
type MockTemporalWorkflowRun_GetID_Call struct {
	*mock.Call
}

// GetID is a helper method to define mock.On call
func (_e *MockTemporalWorkflowRun_Expecter) GetID() *MockTemporalWorkflowRun_GetID_Call {
	return &MockTemporalWorkflowRun_GetID_Call{Call: _e.mock.On("GetID")}
}

func (_c *MockTemporalWorkflowRun_GetID_Call) Run(run func()) *MockTemporalWorkflowRun_GetID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTemporalWorkflowRun_GetID_Call) Return(_a0 string) *MockTemporalWorkflowRun_GetID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemporalWorkflowRun_GetID_Call) RunAndReturn(run func() string) *MockTemporalWorkflowRun_GetID_Call {
	_c.Call.Return(run)
	return _c
}

// GetRunID provides a mock function with given fields:
func (_m *MockTemporalWorkflowRun) GetRunID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockTemporalWorkflowRun_GetRunID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRunID'
type MockTemporalWorkflowRun_GetRunID_Call struct {
	*mock.Call
}

// GetRunID is a helper method to define mock.On call
func (_e *MockTemporalWorkflowRun_Expecter) GetRunID() *MockTemporalWorkflowRun_GetRunID_Call {
	return &MockTemporalWorkflowRun_GetRunID_Call{Call: _e.mock.On("GetRunID")}
}

func (_c *MockTemporalWorkflowRun_GetRunID_Call) Run(run func()) *MockTemporalWorkflowRun_GetRunID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTemporalWorkflowRun_GetRunID_Call) Return(_a0 string) *MockTemporalWorkflowRun_GetRunID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemporalWorkflowRun_GetRunID_Call) RunAndReturn(run func() string) *MockTemporalWorkflowRun_GetRunID_Call {
	_c.Call.Return(run)
	return _c
}

// GetWithOptions provides a mock function with given fields: ctx, valuePtr, options
func (_m *MockTemporalWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, options sdkclient.WorkflowRunGetOptions) error {
	ret := _m.Called(ctx, valuePtr, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, sdkclient.WorkflowRunGetOptions) error); ok {
		r0 = rf(ctx, valuePtr, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemporalWorkflowRun_GetWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWithOptions'
type MockTemporalWorkflowRun_GetWithOptions_Call struct {
	*mock.Call
}

// GetWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - valuePtr interface{}
//   - options sdkclient.WorkflowRunGetOptions
func (_e *MockTemporalWorkflowRun_Expecter) GetWithOptions(ctx interface{}, valuePtr interface{}, options interface{}) *MockTemporalWorkflowRun_GetWithOptions_Call {
	return &MockTemporalWorkflowRun_GetWithOptions_Call{Call: _e.mock.On("GetWithOptions", ctx, valuePtr, options)}
}

func (_c *MockTemporalWorkflowRun_GetWithOptions_Call) Run(run func(ctx context.Context, valuePtr interface{}, options sdkclient.WorkflowRunGetOptions)) *MockTemporalWorkflowRun_GetWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}), args[2].(sdkclient.WorkflowRunGetOptions))
	})
	return _c
}

func (_c *MockTemporalWorkflowRun_GetWithOptions_Call) Return(_a0 error) *MockTemporalWorkflowRun_GetWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemporalWorkflowRun_GetWithOptions_Call) RunAndReturn(run func(context.Context, interface{}, sdkclient.WorkflowRunGetOptions) error) *MockTemporalWorkflowRun_GetWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTemporalWorkflowRun creates a new instance of MockTemporalWorkflowRun. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTemporalWorkflowRun(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTemporalWorkflowRun {
	mock := &MockTemporalWorkflowRun{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
