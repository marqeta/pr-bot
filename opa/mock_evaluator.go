// Code generated by mockery v2.33.0. DO NOT EDIT.

package opa

import (
	context "context"

	input "github.com/marqeta/pr-bot/opa/input"
	mock "github.com/stretchr/testify/mock"

	types "github.com/marqeta/pr-bot/opa/types"
)

// MockEvaluator is an autogenerated mock type for the Evaluator type
type MockEvaluator struct {
	mock.Mock
}

type MockEvaluator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEvaluator) EXPECT() *MockEvaluator_Expecter {
	return &MockEvaluator_Expecter{mock: &_m.Mock}
}

// Evaluate provides a mock function with given fields: ctx, _a1
func (_m *MockEvaluator) Evaluate(ctx context.Context, _a1 input.GHE) (types.Result, error) {
	ret := _m.Called(ctx, _a1)

	var r0 types.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, input.GHE) (types.Result, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, input.GHE) types.Result); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(types.Result)
	}

	if rf, ok := ret.Get(1).(func(context.Context, input.GHE) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEvaluator_Evaluate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Evaluate'
type MockEvaluator_Evaluate_Call struct {
	*mock.Call
}

// Evaluate is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 input.GHE
func (_e *MockEvaluator_Expecter) Evaluate(ctx interface{}, _a1 interface{}) *MockEvaluator_Evaluate_Call {
	return &MockEvaluator_Evaluate_Call{Call: _e.mock.On("Evaluate", ctx, _a1)}
}

func (_c *MockEvaluator_Evaluate_Call) Run(run func(ctx context.Context, _a1 input.GHE)) *MockEvaluator_Evaluate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(input.GHE))
	})
	return _c
}

func (_c *MockEvaluator_Evaluate_Call) Return(_a0 types.Result, _a1 error) *MockEvaluator_Evaluate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEvaluator_Evaluate_Call) RunAndReturn(run func(context.Context, input.GHE) (types.Result, error)) *MockEvaluator_Evaluate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEvaluator creates a new instance of MockEvaluator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEvaluator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEvaluator {
	mock := &MockEvaluator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
