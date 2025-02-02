// Code generated by mockery v2.42.0. DO NOT EDIT.

package mockcontroller

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/alexandremahdhaoui/ipxer/internal/types"
)

// MockIPXE is an autogenerated mock type for the IPXE type
type MockIPXE struct {
	mock.Mock
}

type MockIPXE_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIPXE) EXPECT() *MockIPXE_Expecter {
	return &MockIPXE_Expecter{mock: &_m.Mock}
}

// Boostrap provides a mock function with given fields:
func (_m *MockIPXE) Boostrap() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Boostrap")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockIPXE_Boostrap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Boostrap'
type MockIPXE_Boostrap_Call struct {
	*mock.Call
}

// Boostrap is a helper method to define mock.On call
func (_e *MockIPXE_Expecter) Boostrap() *MockIPXE_Boostrap_Call {
	return &MockIPXE_Boostrap_Call{Call: _e.mock.On("Boostrap")}
}

func (_c *MockIPXE_Boostrap_Call) Run(run func()) *MockIPXE_Boostrap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIPXE_Boostrap_Call) Return(_a0 []byte) *MockIPXE_Boostrap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPXE_Boostrap_Call) RunAndReturn(run func() []byte) *MockIPXE_Boostrap_Call {
	_c.Call.Return(run)
	return _c
}

// FindProfileAndRender provides a mock function with given fields: ctx, selectors
func (_m *MockIPXE) FindProfileAndRender(ctx context.Context, selectors types.IPXESelectors) ([]byte, error) {
	ret := _m.Called(ctx, selectors)

	if len(ret) == 0 {
		panic("no return value specified for FindProfileAndRender")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.IPXESelectors) ([]byte, error)); ok {
		return rf(ctx, selectors)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.IPXESelectors) []byte); ok {
		r0 = rf(ctx, selectors)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.IPXESelectors) error); ok {
		r1 = rf(ctx, selectors)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPXE_FindProfileAndRender_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindProfileAndRender'
type MockIPXE_FindProfileAndRender_Call struct {
	*mock.Call
}

// FindProfileAndRender is a helper method to define mock.On call
//   - ctx context.Context
//   - selectors types.IPXESelectors
func (_e *MockIPXE_Expecter) FindProfileAndRender(ctx interface{}, selectors interface{}) *MockIPXE_FindProfileAndRender_Call {
	return &MockIPXE_FindProfileAndRender_Call{Call: _e.mock.On("FindProfileAndRender", ctx, selectors)}
}

func (_c *MockIPXE_FindProfileAndRender_Call) Run(run func(ctx context.Context, selectors types.IPXESelectors)) *MockIPXE_FindProfileAndRender_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.IPXESelectors))
	})
	return _c
}

func (_c *MockIPXE_FindProfileAndRender_Call) Return(_a0 []byte, _a1 error) *MockIPXE_FindProfileAndRender_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPXE_FindProfileAndRender_Call) RunAndReturn(run func(context.Context, types.IPXESelectors) ([]byte, error)) *MockIPXE_FindProfileAndRender_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIPXE creates a new instance of MockIPXE. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIPXE(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIPXE {
	mock := &MockIPXE{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
