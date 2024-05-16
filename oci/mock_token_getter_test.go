// Code generated by mockery v2.33.0. DO NOT EDIT.

package oci

import (
	context "context"

	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	mock "github.com/stretchr/testify/mock"
)

// MockTokenGetter is an autogenerated mock type for the TokenGetter type
type MockTokenGetter struct {
	mock.Mock
}

type MockTokenGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTokenGetter) EXPECT() *MockTokenGetter_Expecter {
	return &MockTokenGetter_Expecter{mock: &_m.Mock}
}

// GetAuthorizationToken provides a mock function with given fields: ctx, input, optFns
func (_m *MockTokenGetter) GetAuthorizationToken(ctx context.Context, input *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ecr.GetAuthorizationTokenOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)); ok {
		return rf(ctx, input, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) *ecr.GetAuthorizationTokenOutput); ok {
		r0 = rf(ctx, input, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecr.GetAuthorizationTokenOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) error); ok {
		r1 = rf(ctx, input, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTokenGetter_GetAuthorizationToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAuthorizationToken'
type MockTokenGetter_GetAuthorizationToken_Call struct {
	*mock.Call
}

// GetAuthorizationToken is a helper method to define mock.On call
//   - ctx context.Context
//   - input *ecr.GetAuthorizationTokenInput
//   - optFns ...func(*ecr.Options)
func (_e *MockTokenGetter_Expecter) GetAuthorizationToken(ctx interface{}, input interface{}, optFns ...interface{}) *MockTokenGetter_GetAuthorizationToken_Call {
	return &MockTokenGetter_GetAuthorizationToken_Call{Call: _e.mock.On("GetAuthorizationToken",
		append([]interface{}{ctx, input}, optFns...)...)}
}

func (_c *MockTokenGetter_GetAuthorizationToken_Call) Run(run func(ctx context.Context, input *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options))) *MockTokenGetter_GetAuthorizationToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*ecr.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*ecr.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*ecr.GetAuthorizationTokenInput), variadicArgs...)
	})
	return _c
}

func (_c *MockTokenGetter_GetAuthorizationToken_Call) Return(_a0 *ecr.GetAuthorizationTokenOutput, _a1 error) *MockTokenGetter_GetAuthorizationToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTokenGetter_GetAuthorizationToken_Call) RunAndReturn(run func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)) *MockTokenGetter_GetAuthorizationToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTokenGetter creates a new instance of MockTokenGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTokenGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTokenGetter {
	mock := &MockTokenGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}