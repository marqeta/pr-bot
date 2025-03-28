// Code generated by mockery v2.49.0. DO NOT EDIT.

package input

import (
	context "context"

	github "github.com/google/go-github/v50/github"
	datastore "github.com/marqeta/pr-bot/datastore"

	mock "github.com/stretchr/testify/mock"
)

// MockAdapter is an autogenerated mock type for the Adapter type
type MockAdapter struct {
	mock.Mock
}

type MockAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAdapter) EXPECT() *MockAdapter_Expecter {
	return &MockAdapter_Expecter{mock: &_m.Mock}
}

// MetadataToGHE provides a mock function with given fields: ctx, metadata
func (_m *MockAdapter) MetadataToGHE(ctx context.Context, metadata *datastore.Metadata) (GHE, error) {
	ret := _m.Called(ctx, metadata)

	if len(ret) == 0 {
		panic("no return value specified for MetadataToGHE")
	}

	var r0 GHE
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *datastore.Metadata) (GHE, error)); ok {
		return rf(ctx, metadata)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *datastore.Metadata) GHE); ok {
		r0 = rf(ctx, metadata)
	} else {
		r0 = ret.Get(0).(GHE)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *datastore.Metadata) error); ok {
		r1 = rf(ctx, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapter_MetadataToGHE_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MetadataToGHE'
type MockAdapter_MetadataToGHE_Call struct {
	*mock.Call
}

// MetadataToGHE is a helper method to define mock.On call
//   - ctx context.Context
//   - metadata *datastore.Metadata
func (_e *MockAdapter_Expecter) MetadataToGHE(ctx interface{}, metadata interface{}) *MockAdapter_MetadataToGHE_Call {
	return &MockAdapter_MetadataToGHE_Call{Call: _e.mock.On("MetadataToGHE", ctx, metadata)}
}

func (_c *MockAdapter_MetadataToGHE_Call) Run(run func(ctx context.Context, metadata *datastore.Metadata)) *MockAdapter_MetadataToGHE_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*datastore.Metadata))
	})
	return _c
}

func (_c *MockAdapter_MetadataToGHE_Call) Return(_a0 GHE, _a1 error) *MockAdapter_MetadataToGHE_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapter_MetadataToGHE_Call) RunAndReturn(run func(context.Context, *datastore.Metadata) (GHE, error)) *MockAdapter_MetadataToGHE_Call {
	_c.Call.Return(run)
	return _c
}

// PREventToGHE provides a mock function with given fields: ctx, event
func (_m *MockAdapter) PREventToGHE(ctx context.Context, event *github.PullRequestEvent) (GHE, error) {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for PREventToGHE")
	}

	var r0 GHE
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *github.PullRequestEvent) (GHE, error)); ok {
		return rf(ctx, event)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *github.PullRequestEvent) GHE); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(GHE)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *github.PullRequestEvent) error); ok {
		r1 = rf(ctx, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapter_PREventToGHE_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PREventToGHE'
type MockAdapter_PREventToGHE_Call struct {
	*mock.Call
}

// PREventToGHE is a helper method to define mock.On call
//   - ctx context.Context
//   - event *github.PullRequestEvent
func (_e *MockAdapter_Expecter) PREventToGHE(ctx interface{}, event interface{}) *MockAdapter_PREventToGHE_Call {
	return &MockAdapter_PREventToGHE_Call{Call: _e.mock.On("PREventToGHE", ctx, event)}
}

func (_c *MockAdapter_PREventToGHE_Call) Run(run func(ctx context.Context, event *github.PullRequestEvent)) *MockAdapter_PREventToGHE_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*github.PullRequestEvent))
	})
	return _c
}

func (_c *MockAdapter_PREventToGHE_Call) Return(_a0 GHE, _a1 error) *MockAdapter_PREventToGHE_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapter_PREventToGHE_Call) RunAndReturn(run func(context.Context, *github.PullRequestEvent) (GHE, error)) *MockAdapter_PREventToGHE_Call {
	_c.Call.Return(run)
	return _c
}

// PRReviewEventToGHE provides a mock function with given fields: ctx, event
func (_m *MockAdapter) PRReviewEventToGHE(ctx context.Context, event *github.PullRequestReviewEvent) (GHE, error) {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for PRReviewEventToGHE")
	}

	var r0 GHE
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *github.PullRequestReviewEvent) (GHE, error)); ok {
		return rf(ctx, event)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *github.PullRequestReviewEvent) GHE); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(GHE)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *github.PullRequestReviewEvent) error); ok {
		r1 = rf(ctx, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapter_PRReviewEventToGHE_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PRReviewEventToGHE'
type MockAdapter_PRReviewEventToGHE_Call struct {
	*mock.Call
}

// PRReviewEventToGHE is a helper method to define mock.On call
//   - ctx context.Context
//   - event *github.PullRequestReviewEvent
func (_e *MockAdapter_Expecter) PRReviewEventToGHE(ctx interface{}, event interface{}) *MockAdapter_PRReviewEventToGHE_Call {
	return &MockAdapter_PRReviewEventToGHE_Call{Call: _e.mock.On("PRReviewEventToGHE", ctx, event)}
}

func (_c *MockAdapter_PRReviewEventToGHE_Call) Run(run func(ctx context.Context, event *github.PullRequestReviewEvent)) *MockAdapter_PRReviewEventToGHE_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*github.PullRequestReviewEvent))
	})
	return _c
}

func (_c *MockAdapter_PRReviewEventToGHE_Call) Return(_a0 GHE, _a1 error) *MockAdapter_PRReviewEventToGHE_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapter_PRReviewEventToGHE_Call) RunAndReturn(run func(context.Context, *github.PullRequestReviewEvent) (GHE, error)) *MockAdapter_PRReviewEventToGHE_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAdapter creates a new instance of MockAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAdapter {
	mock := &MockAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
