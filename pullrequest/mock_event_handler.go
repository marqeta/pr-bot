// Code generated by mockery v2.49.0. DO NOT EDIT.

package pullrequest

import (
	context "context"

	github "github.com/google/go-github/v50/github"
	datastore "github.com/marqeta/pr-bot/datastore"

	id "github.com/marqeta/pr-bot/id"

	input "github.com/marqeta/pr-bot/opa/input"

	mock "github.com/stretchr/testify/mock"
)

// MockEventHandler is an autogenerated mock type for the EventHandler type
type MockEventHandler struct {
	mock.Mock
}

type MockEventHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventHandler) EXPECT() *MockEventHandler_Expecter {
	return &MockEventHandler_Expecter{mock: &_m.Mock}
}

// EvalAndReview provides a mock function with given fields: ctx, _a1, ghe
func (_m *MockEventHandler) EvalAndReview(ctx context.Context, _a1 id.PR, ghe input.GHE) error {
	ret := _m.Called(ctx, _a1, ghe)

	if len(ret) == 0 {
		panic("no return value specified for EvalAndReview")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, input.GHE) error); ok {
		r0 = rf(ctx, _a1, ghe)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventHandler_EvalAndReview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EvalAndReview'
type MockEventHandler_EvalAndReview_Call struct {
	*mock.Call
}

// EvalAndReview is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - ghe input.GHE
func (_e *MockEventHandler_Expecter) EvalAndReview(ctx interface{}, _a1 interface{}, ghe interface{}) *MockEventHandler_EvalAndReview_Call {
	return &MockEventHandler_EvalAndReview_Call{Call: _e.mock.On("EvalAndReview", ctx, _a1, ghe)}
}

func (_c *MockEventHandler_EvalAndReview_Call) Run(run func(ctx context.Context, _a1 id.PR, ghe input.GHE)) *MockEventHandler_EvalAndReview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(input.GHE))
	})
	return _c
}

func (_c *MockEventHandler_EvalAndReview_Call) Return(_a0 error) *MockEventHandler_EvalAndReview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_EvalAndReview_Call) RunAndReturn(run func(context.Context, id.PR, input.GHE) error) *MockEventHandler_EvalAndReview_Call {
	_c.Call.Return(run)
	return _c
}

// EvalAndReviewDataEvent provides a mock function with given fields: ctx, metadata
func (_m *MockEventHandler) EvalAndReviewDataEvent(ctx context.Context, metadata *datastore.Metadata) error {
	ret := _m.Called(ctx, metadata)

	if len(ret) == 0 {
		panic("no return value specified for EvalAndReviewDataEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *datastore.Metadata) error); ok {
		r0 = rf(ctx, metadata)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventHandler_EvalAndReviewDataEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EvalAndReviewDataEvent'
type MockEventHandler_EvalAndReviewDataEvent_Call struct {
	*mock.Call
}

// EvalAndReviewDataEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - metadata *datastore.Metadata
func (_e *MockEventHandler_Expecter) EvalAndReviewDataEvent(ctx interface{}, metadata interface{}) *MockEventHandler_EvalAndReviewDataEvent_Call {
	return &MockEventHandler_EvalAndReviewDataEvent_Call{Call: _e.mock.On("EvalAndReviewDataEvent", ctx, metadata)}
}

func (_c *MockEventHandler_EvalAndReviewDataEvent_Call) Run(run func(ctx context.Context, metadata *datastore.Metadata)) *MockEventHandler_EvalAndReviewDataEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*datastore.Metadata))
	})
	return _c
}

func (_c *MockEventHandler_EvalAndReviewDataEvent_Call) Return(_a0 error) *MockEventHandler_EvalAndReviewDataEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_EvalAndReviewDataEvent_Call) RunAndReturn(run func(context.Context, *datastore.Metadata) error) *MockEventHandler_EvalAndReviewDataEvent_Call {
	_c.Call.Return(run)
	return _c
}

// EvalAndReviewPREvent provides a mock function with given fields: ctx, _a1, event
func (_m *MockEventHandler) EvalAndReviewPREvent(ctx context.Context, _a1 id.PR, event *github.PullRequestEvent) error {
	ret := _m.Called(ctx, _a1, event)

	if len(ret) == 0 {
		panic("no return value specified for EvalAndReviewPREvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, *github.PullRequestEvent) error); ok {
		r0 = rf(ctx, _a1, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventHandler_EvalAndReviewPREvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EvalAndReviewPREvent'
type MockEventHandler_EvalAndReviewPREvent_Call struct {
	*mock.Call
}

// EvalAndReviewPREvent is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - event *github.PullRequestEvent
func (_e *MockEventHandler_Expecter) EvalAndReviewPREvent(ctx interface{}, _a1 interface{}, event interface{}) *MockEventHandler_EvalAndReviewPREvent_Call {
	return &MockEventHandler_EvalAndReviewPREvent_Call{Call: _e.mock.On("EvalAndReviewPREvent", ctx, _a1, event)}
}

func (_c *MockEventHandler_EvalAndReviewPREvent_Call) Run(run func(ctx context.Context, _a1 id.PR, event *github.PullRequestEvent)) *MockEventHandler_EvalAndReviewPREvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(*github.PullRequestEvent))
	})
	return _c
}

func (_c *MockEventHandler_EvalAndReviewPREvent_Call) Return(_a0 error) *MockEventHandler_EvalAndReviewPREvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_EvalAndReviewPREvent_Call) RunAndReturn(run func(context.Context, id.PR, *github.PullRequestEvent) error) *MockEventHandler_EvalAndReviewPREvent_Call {
	_c.Call.Return(run)
	return _c
}

// EvalAndReviewPRReviewEvent provides a mock function with given fields: ctx, _a1, event
func (_m *MockEventHandler) EvalAndReviewPRReviewEvent(ctx context.Context, _a1 id.PR, event *github.PullRequestReviewEvent) error {
	ret := _m.Called(ctx, _a1, event)

	if len(ret) == 0 {
		panic("no return value specified for EvalAndReviewPRReviewEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, *github.PullRequestReviewEvent) error); ok {
		r0 = rf(ctx, _a1, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventHandler_EvalAndReviewPRReviewEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EvalAndReviewPRReviewEvent'
type MockEventHandler_EvalAndReviewPRReviewEvent_Call struct {
	*mock.Call
}

// EvalAndReviewPRReviewEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - event *github.PullRequestReviewEvent
func (_e *MockEventHandler_Expecter) EvalAndReviewPRReviewEvent(ctx interface{}, _a1 interface{}, event interface{}) *MockEventHandler_EvalAndReviewPRReviewEvent_Call {
	return &MockEventHandler_EvalAndReviewPRReviewEvent_Call{Call: _e.mock.On("EvalAndReviewPRReviewEvent", ctx, _a1, event)}
}

func (_c *MockEventHandler_EvalAndReviewPRReviewEvent_Call) Run(run func(ctx context.Context, _a1 id.PR, event *github.PullRequestReviewEvent)) *MockEventHandler_EvalAndReviewPRReviewEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(*github.PullRequestReviewEvent))
	})
	return _c
}

func (_c *MockEventHandler_EvalAndReviewPRReviewEvent_Call) Return(_a0 error) *MockEventHandler_EvalAndReviewPRReviewEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_EvalAndReviewPRReviewEvent_Call) RunAndReturn(run func(context.Context, id.PR, *github.PullRequestReviewEvent) error) *MockEventHandler_EvalAndReviewPRReviewEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventHandler creates a new instance of MockEventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventHandler {
	mock := &MockEventHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
