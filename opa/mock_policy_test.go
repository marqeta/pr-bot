// Code generated by mockery v2.33.0. DO NOT EDIT.

package opa

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	input "github.com/marqeta/pr-bot/opa/input"

	types "github.com/marqeta/pr-bot/opa/types"
)

// MockPolicy is an autogenerated mock type for the Policy type
type MockPolicy struct {
	mock.Mock
}

type MockPolicy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPolicy) EXPECT() *MockPolicy_Expecter {
	return &MockPolicy_Expecter{mock: &_m.Mock}
}

// Evaluate provides a mock function with given fields: ctx, module, _a2
func (_m *MockPolicy) Evaluate(ctx context.Context, module string, _a2 *input.Model) (types.Result, error) {
	ret := _m.Called(ctx, module, _a2)

	var r0 types.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *input.Model) (types.Result, error)); ok {
		return rf(ctx, module, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *input.Model) types.Result); ok {
		r0 = rf(ctx, module, _a2)
	} else {
		r0 = ret.Get(0).(types.Result)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *input.Model) error); ok {
		r1 = rf(ctx, module, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPolicy_Evaluate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Evaluate'
type MockPolicy_Evaluate_Call struct {
	*mock.Call
}

// Evaluate is a helper method to define mock.On call
//   - ctx context.Context
//   - module string
//   - _a2 *input.Model
func (_e *MockPolicy_Expecter) Evaluate(ctx interface{}, module interface{}, _a2 interface{}) *MockPolicy_Evaluate_Call {
	return &MockPolicy_Evaluate_Call{Call: _e.mock.On("Evaluate", ctx, module, _a2)}
}

func (_c *MockPolicy_Evaluate_Call) Run(run func(ctx context.Context, module string, _a2 *input.Model)) *MockPolicy_Evaluate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*input.Model))
	})
	return _c
}

func (_c *MockPolicy_Evaluate_Call) Return(_a0 types.Result, _a1 error) *MockPolicy_Evaluate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicy_Evaluate_Call) RunAndReturn(run func(context.Context, string, *input.Model) (types.Result, error)) *MockPolicy_Evaluate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPolicy creates a new instance of MockPolicy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPolicy(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPolicy {
	mock := &MockPolicy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}