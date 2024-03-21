// Code generated by mockery v2.42.0. DO NOT EDIT.

package mockadapters

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	types "github.com/alexandremahdhaoui/ipxer/internal/types"
)

// MockWebhookResolver is an autogenerated mock type for the WebhookResolver type
type MockWebhookResolver struct {
	mock.Mock
}

type MockWebhookResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWebhookResolver) EXPECT() *MockWebhookResolver_Expecter {
	return &MockWebhookResolver_Expecter{mock: &_m.Mock}
}

// Resolve provides a mock function with given fields: ctx, c
func (_m *MockWebhookResolver) Resolve(ctx context.Context, c types.Content) ([]byte, error) {
	ret := _m.Called(ctx, c)

	if len(ret) == 0 {
		panic("no return value specified for Resolve")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Content) ([]byte, error)); ok {
		return rf(ctx, c)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Content) []byte); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Content) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWebhookResolver_Resolve_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Resolve'
type MockWebhookResolver_Resolve_Call struct {
	*mock.Call
}

// Resolve is a helper method to define mock.On call
//   - ctx context.Context
//   - c types.Content
func (_e *MockWebhookResolver_Expecter) Resolve(ctx interface{}, c interface{}) *MockWebhookResolver_Resolve_Call {
	return &MockWebhookResolver_Resolve_Call{Call: _e.mock.On("Resolve", ctx, c)}
}

func (_c *MockWebhookResolver_Resolve_Call) Run(run func(ctx context.Context, c types.Content)) *MockWebhookResolver_Resolve_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Content))
	})
	return _c
}

func (_c *MockWebhookResolver_Resolve_Call) Return(_a0 []byte, _a1 error) *MockWebhookResolver_Resolve_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockWebhookResolver_Resolve_Call) RunAndReturn(run func(context.Context, types.Content) ([]byte, error)) *MockWebhookResolver_Resolve_Call {
	_c.Call.Return(run)
	return _c
}

// ResolveRequest provides a mock function with given fields: ctx, req, config
func (_m *MockWebhookResolver) ResolveRequest(ctx context.Context, req *http.Request, config types.WebhookConfig) ([][]byte, error) {
	ret := _m.Called(ctx, req, config)

	if len(ret) == 0 {
		panic("no return value specified for ResolveRequest")
	}

	var r0 [][]byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request, types.WebhookConfig) ([][]byte, error)); ok {
		return rf(ctx, req, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request, types.WebhookConfig) [][]byte); ok {
		r0 = rf(ctx, req, config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *http.Request, types.WebhookConfig) error); ok {
		r1 = rf(ctx, req, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWebhookResolver_ResolveRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResolveRequest'
type MockWebhookResolver_ResolveRequest_Call struct {
	*mock.Call
}

// ResolveRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - req *http.Request
//   - config types.WebhookConfig
func (_e *MockWebhookResolver_Expecter) ResolveRequest(ctx interface{}, req interface{}, config interface{}) *MockWebhookResolver_ResolveRequest_Call {
	return &MockWebhookResolver_ResolveRequest_Call{Call: _e.mock.On("ResolveRequest", ctx, req, config)}
}

func (_c *MockWebhookResolver_ResolveRequest_Call) Run(run func(ctx context.Context, req *http.Request, config types.WebhookConfig)) *MockWebhookResolver_ResolveRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*http.Request), args[2].(types.WebhookConfig))
	})
	return _c
}

func (_c *MockWebhookResolver_ResolveRequest_Call) Return(_a0 [][]byte, _a1 error) *MockWebhookResolver_ResolveRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockWebhookResolver_ResolveRequest_Call) RunAndReturn(run func(context.Context, *http.Request, types.WebhookConfig) ([][]byte, error)) *MockWebhookResolver_ResolveRequest_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWebhookResolver creates a new instance of MockWebhookResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWebhookResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWebhookResolver {
	mock := &MockWebhookResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
