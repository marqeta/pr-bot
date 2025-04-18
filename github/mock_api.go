// Code generated by mockery v2.49.0. DO NOT EDIT.

package github

import (
	context "context"

	errors "github.com/marqeta/pr-bot/errors"
	githubv4 "github.com/shurcooL/githubv4"

	id "github.com/marqeta/pr-bot/id"

	mock "github.com/stretchr/testify/mock"

	v50github "github.com/google/go-github/v50/github"
)

// MockAPI is an autogenerated mock type for the API type
type MockAPI struct {
	mock.Mock
}

type MockAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAPI) EXPECT() *MockAPI_Expecter {
	return &MockAPI_Expecter{mock: &_m.Mock}
}

// AddReview provides a mock function with given fields: ctx, _a1, summary, event
func (_m *MockAPI) AddReview(ctx context.Context, _a1 id.PR, summary string, event string) error {
	ret := _m.Called(ctx, _a1, summary, event)

	if len(ret) == 0 {
		panic("no return value specified for AddReview")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string, string) error); ok {
		r0 = rf(ctx, _a1, summary, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAPI_AddReview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddReview'
type MockAPI_AddReview_Call struct {
	*mock.Call
}

// AddReview is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - summary string
//   - event string
func (_e *MockAPI_Expecter) AddReview(ctx interface{}, _a1 interface{}, summary interface{}, event interface{}) *MockAPI_AddReview_Call {
	return &MockAPI_AddReview_Call{Call: _e.mock.On("AddReview", ctx, _a1, summary, event)}
}

func (_c *MockAPI_AddReview_Call) Run(run func(ctx context.Context, _a1 id.PR, summary string, event string)) *MockAPI_AddReview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockAPI_AddReview_Call) Return(_a0 error) *MockAPI_AddReview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAPI_AddReview_Call) RunAndReturn(run func(context.Context, id.PR, string, string) error) *MockAPI_AddReview_Call {
	_c.Call.Return(run)
	return _c
}

// EnableAutoMerge provides a mock function with given fields: ctx, _a1, method
func (_m *MockAPI) EnableAutoMerge(ctx context.Context, _a1 id.PR, method githubv4.PullRequestMergeMethod) error {
	ret := _m.Called(ctx, _a1, method)

	if len(ret) == 0 {
		panic("no return value specified for EnableAutoMerge")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, githubv4.PullRequestMergeMethod) error); ok {
		r0 = rf(ctx, _a1, method)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAPI_EnableAutoMerge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EnableAutoMerge'
type MockAPI_EnableAutoMerge_Call struct {
	*mock.Call
}

// EnableAutoMerge is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - method githubv4.PullRequestMergeMethod
func (_e *MockAPI_Expecter) EnableAutoMerge(ctx interface{}, _a1 interface{}, method interface{}) *MockAPI_EnableAutoMerge_Call {
	return &MockAPI_EnableAutoMerge_Call{Call: _e.mock.On("EnableAutoMerge", ctx, _a1, method)}
}

func (_c *MockAPI_EnableAutoMerge_Call) Run(run func(ctx context.Context, _a1 id.PR, method githubv4.PullRequestMergeMethod)) *MockAPI_EnableAutoMerge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(githubv4.PullRequestMergeMethod))
	})
	return _c
}

func (_c *MockAPI_EnableAutoMerge_Call) Return(_a0 error) *MockAPI_EnableAutoMerge_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAPI_EnableAutoMerge_Call) RunAndReturn(run func(context.Context, id.PR, githubv4.PullRequestMergeMethod) error) *MockAPI_EnableAutoMerge_Call {
	_c.Call.Return(run)
	return _c
}

// GetBranchProtection provides a mock function with given fields: ctx, _a1, branch
func (_m *MockAPI) GetBranchProtection(ctx context.Context, _a1 id.PR, branch string) (*v50github.Protection, error) {
	ret := _m.Called(ctx, _a1, branch)

	if len(ret) == 0 {
		panic("no return value specified for GetBranchProtection")
	}

	var r0 *v50github.Protection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) (*v50github.Protection, error)); ok {
		return rf(ctx, _a1, branch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) *v50github.Protection); ok {
		r0 = rf(ctx, _a1, branch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v50github.Protection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR, string) error); ok {
		r1 = rf(ctx, _a1, branch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_GetBranchProtection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBranchProtection'
type MockAPI_GetBranchProtection_Call struct {
	*mock.Call
}

// GetBranchProtection is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - branch string
func (_e *MockAPI_Expecter) GetBranchProtection(ctx interface{}, _a1 interface{}, branch interface{}) *MockAPI_GetBranchProtection_Call {
	return &MockAPI_GetBranchProtection_Call{Call: _e.mock.On("GetBranchProtection", ctx, _a1, branch)}
}

func (_c *MockAPI_GetBranchProtection_Call) Run(run func(ctx context.Context, _a1 id.PR, branch string)) *MockAPI_GetBranchProtection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(string))
	})
	return _c
}

func (_c *MockAPI_GetBranchProtection_Call) Return(_a0 *v50github.Protection, _a1 error) *MockAPI_GetBranchProtection_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_GetBranchProtection_Call) RunAndReturn(run func(context.Context, id.PR, string) (*v50github.Protection, error)) *MockAPI_GetBranchProtection_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganization provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) GetOrganization(ctx context.Context, _a1 id.PR) (*v50github.Organization, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganization")
	}

	var r0 *v50github.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) (*v50github.Organization, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) *v50github.Organization); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v50github.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_GetOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganization'
type MockAPI_GetOrganization_Call struct {
	*mock.Call
}

// GetOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) GetOrganization(ctx interface{}, _a1 interface{}) *MockAPI_GetOrganization_Call {
	return &MockAPI_GetOrganization_Call{Call: _e.mock.On("GetOrganization", ctx, _a1)}
}

func (_c *MockAPI_GetOrganization_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_GetOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_GetOrganization_Call) Return(_a0 *v50github.Organization, _a1 error) *MockAPI_GetOrganization_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_GetOrganization_Call) RunAndReturn(run func(context.Context, id.PR) (*v50github.Organization, error)) *MockAPI_GetOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// GetPullRequest provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) GetPullRequest(ctx context.Context, _a1 id.PR) (*v50github.PullRequest, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetPullRequest")
	}

	var r0 *v50github.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) (*v50github.PullRequest, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) *v50github.PullRequest); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v50github.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_GetPullRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPullRequest'
type MockAPI_GetPullRequest_Call struct {
	*mock.Call
}

// GetPullRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) GetPullRequest(ctx interface{}, _a1 interface{}) *MockAPI_GetPullRequest_Call {
	return &MockAPI_GetPullRequest_Call{Call: _e.mock.On("GetPullRequest", ctx, _a1)}
}

func (_c *MockAPI_GetPullRequest_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_GetPullRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_GetPullRequest_Call) Return(_a0 *v50github.PullRequest, _a1 error) *MockAPI_GetPullRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_GetPullRequest_Call) RunAndReturn(run func(context.Context, id.PR) (*v50github.PullRequest, error)) *MockAPI_GetPullRequest_Call {
	_c.Call.Return(run)
	return _c
}

// GetRepository provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) GetRepository(ctx context.Context, _a1 id.PR) (*v50github.Repository, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetRepository")
	}

	var r0 *v50github.Repository
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) (*v50github.Repository, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) *v50github.Repository); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v50github.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_GetRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRepository'
type MockAPI_GetRepository_Call struct {
	*mock.Call
}

// GetRepository is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) GetRepository(ctx interface{}, _a1 interface{}) *MockAPI_GetRepository_Call {
	return &MockAPI_GetRepository_Call{Call: _e.mock.On("GetRepository", ctx, _a1)}
}

func (_c *MockAPI_GetRepository_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_GetRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_GetRepository_Call) Return(_a0 *v50github.Repository, _a1 error) *MockAPI_GetRepository_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_GetRepository_Call) RunAndReturn(run func(context.Context, id.PR) (*v50github.Repository, error)) *MockAPI_GetRepository_Call {
	_c.Call.Return(run)
	return _c
}

// IssueComment provides a mock function with given fields: ctx, _a1, comment
func (_m *MockAPI) IssueComment(ctx context.Context, _a1 id.PR, comment string) error {
	ret := _m.Called(ctx, _a1, comment)

	if len(ret) == 0 {
		panic("no return value specified for IssueComment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) error); ok {
		r0 = rf(ctx, _a1, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAPI_IssueComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IssueComment'
type MockAPI_IssueComment_Call struct {
	*mock.Call
}

// IssueComment is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - comment string
func (_e *MockAPI_Expecter) IssueComment(ctx interface{}, _a1 interface{}, comment interface{}) *MockAPI_IssueComment_Call {
	return &MockAPI_IssueComment_Call{Call: _e.mock.On("IssueComment", ctx, _a1, comment)}
}

func (_c *MockAPI_IssueComment_Call) Run(run func(ctx context.Context, _a1 id.PR, comment string)) *MockAPI_IssueComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(string))
	})
	return _c
}

func (_c *MockAPI_IssueComment_Call) Return(_a0 error) *MockAPI_IssueComment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAPI_IssueComment_Call) RunAndReturn(run func(context.Context, id.PR, string) error) *MockAPI_IssueComment_Call {
	_c.Call.Return(run)
	return _c
}

// IssueCommentForError provides a mock function with given fields: ctx, _a1, err
func (_m *MockAPI) IssueCommentForError(ctx context.Context, _a1 id.PR, err errors.APIError) error {
	ret := _m.Called(ctx, _a1, err)

	if len(ret) == 0 {
		panic("no return value specified for IssueCommentForError")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, errors.APIError) error); ok {
		r0 = rf(ctx, _a1, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAPI_IssueCommentForError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IssueCommentForError'
type MockAPI_IssueCommentForError_Call struct {
	*mock.Call
}

// IssueCommentForError is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - err errors.APIError
func (_e *MockAPI_Expecter) IssueCommentForError(ctx interface{}, _a1 interface{}, err interface{}) *MockAPI_IssueCommentForError_Call {
	return &MockAPI_IssueCommentForError_Call{Call: _e.mock.On("IssueCommentForError", ctx, _a1, err)}
}

func (_c *MockAPI_IssueCommentForError_Call) Run(run func(ctx context.Context, _a1 id.PR, err errors.APIError)) *MockAPI_IssueCommentForError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(errors.APIError))
	})
	return _c
}

func (_c *MockAPI_IssueCommentForError_Call) Return(_a0 error) *MockAPI_IssueCommentForError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAPI_IssueCommentForError_Call) RunAndReturn(run func(context.Context, id.PR, errors.APIError) error) *MockAPI_IssueCommentForError_Call {
	_c.Call.Return(run)
	return _c
}

// ListAllTopics provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) ListAllTopics(ctx context.Context, _a1 id.PR) ([]string, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAllTopics")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) ([]string, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) []string); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListAllTopics_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAllTopics'
type MockAPI_ListAllTopics_Call struct {
	*mock.Call
}

// ListAllTopics is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) ListAllTopics(ctx interface{}, _a1 interface{}) *MockAPI_ListAllTopics_Call {
	return &MockAPI_ListAllTopics_Call{Call: _e.mock.On("ListAllTopics", ctx, _a1)}
}

func (_c *MockAPI_ListAllTopics_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_ListAllTopics_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_ListAllTopics_Call) Return(_a0 []string, _a1 error) *MockAPI_ListAllTopics_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListAllTopics_Call) RunAndReturn(run func(context.Context, id.PR) ([]string, error)) *MockAPI_ListAllTopics_Call {
	_c.Call.Return(run)
	return _c
}

// ListFilesChangedInPR provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) ListFilesChangedInPR(ctx context.Context, _a1 id.PR) ([]*v50github.CommitFile, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListFilesChangedInPR")
	}

	var r0 []*v50github.CommitFile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) ([]*v50github.CommitFile, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) []*v50github.CommitFile); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v50github.CommitFile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListFilesChangedInPR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListFilesChangedInPR'
type MockAPI_ListFilesChangedInPR_Call struct {
	*mock.Call
}

// ListFilesChangedInPR is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) ListFilesChangedInPR(ctx interface{}, _a1 interface{}) *MockAPI_ListFilesChangedInPR_Call {
	return &MockAPI_ListFilesChangedInPR_Call{Call: _e.mock.On("ListFilesChangedInPR", ctx, _a1)}
}

func (_c *MockAPI_ListFilesChangedInPR_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_ListFilesChangedInPR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_ListFilesChangedInPR_Call) Return(_a0 []*v50github.CommitFile, _a1 error) *MockAPI_ListFilesChangedInPR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListFilesChangedInPR_Call) RunAndReturn(run func(context.Context, id.PR) ([]*v50github.CommitFile, error)) *MockAPI_ListFilesChangedInPR_Call {
	_c.Call.Return(run)
	return _c
}

// ListFilesInRootDir provides a mock function with given fields: ctx, _a1, branch
func (_m *MockAPI) ListFilesInRootDir(ctx context.Context, _a1 id.PR, branch string) ([]string, error) {
	ret := _m.Called(ctx, _a1, branch)

	if len(ret) == 0 {
		panic("no return value specified for ListFilesInRootDir")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) ([]string, error)); ok {
		return rf(ctx, _a1, branch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) []string); ok {
		r0 = rf(ctx, _a1, branch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR, string) error); ok {
		r1 = rf(ctx, _a1, branch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListFilesInRootDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListFilesInRootDir'
type MockAPI_ListFilesInRootDir_Call struct {
	*mock.Call
}

// ListFilesInRootDir is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - branch string
func (_e *MockAPI_Expecter) ListFilesInRootDir(ctx interface{}, _a1 interface{}, branch interface{}) *MockAPI_ListFilesInRootDir_Call {
	return &MockAPI_ListFilesInRootDir_Call{Call: _e.mock.On("ListFilesInRootDir", ctx, _a1, branch)}
}

func (_c *MockAPI_ListFilesInRootDir_Call) Run(run func(ctx context.Context, _a1 id.PR, branch string)) *MockAPI_ListFilesInRootDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(string))
	})
	return _c
}

func (_c *MockAPI_ListFilesInRootDir_Call) Return(_a0 []string, _a1 error) *MockAPI_ListFilesInRootDir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListFilesInRootDir_Call) RunAndReturn(run func(context.Context, id.PR, string) ([]string, error)) *MockAPI_ListFilesInRootDir_Call {
	_c.Call.Return(run)
	return _c
}

// ListNamesOfFilesChangedInPR provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) ListNamesOfFilesChangedInPR(ctx context.Context, _a1 id.PR) ([]string, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListNamesOfFilesChangedInPR")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) ([]string, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) []string); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListNamesOfFilesChangedInPR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNamesOfFilesChangedInPR'
type MockAPI_ListNamesOfFilesChangedInPR_Call struct {
	*mock.Call
}

// ListNamesOfFilesChangedInPR is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) ListNamesOfFilesChangedInPR(ctx interface{}, _a1 interface{}) *MockAPI_ListNamesOfFilesChangedInPR_Call {
	return &MockAPI_ListNamesOfFilesChangedInPR_Call{Call: _e.mock.On("ListNamesOfFilesChangedInPR", ctx, _a1)}
}

func (_c *MockAPI_ListNamesOfFilesChangedInPR_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_ListNamesOfFilesChangedInPR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_ListNamesOfFilesChangedInPR_Call) Return(_a0 []string, _a1 error) *MockAPI_ListNamesOfFilesChangedInPR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListNamesOfFilesChangedInPR_Call) RunAndReturn(run func(context.Context, id.PR) ([]string, error)) *MockAPI_ListNamesOfFilesChangedInPR_Call {
	_c.Call.Return(run)
	return _c
}

// ListRequiredStatusChecks provides a mock function with given fields: ctx, _a1, branch
func (_m *MockAPI) ListRequiredStatusChecks(ctx context.Context, _a1 id.PR, branch string) ([]string, error) {
	ret := _m.Called(ctx, _a1, branch)

	if len(ret) == 0 {
		panic("no return value specified for ListRequiredStatusChecks")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) ([]string, error)); ok {
		return rf(ctx, _a1, branch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR, string) []string); ok {
		r0 = rf(ctx, _a1, branch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR, string) error); ok {
		r1 = rf(ctx, _a1, branch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListRequiredStatusChecks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListRequiredStatusChecks'
type MockAPI_ListRequiredStatusChecks_Call struct {
	*mock.Call
}

// ListRequiredStatusChecks is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
//   - branch string
func (_e *MockAPI_Expecter) ListRequiredStatusChecks(ctx interface{}, _a1 interface{}, branch interface{}) *MockAPI_ListRequiredStatusChecks_Call {
	return &MockAPI_ListRequiredStatusChecks_Call{Call: _e.mock.On("ListRequiredStatusChecks", ctx, _a1, branch)}
}

func (_c *MockAPI_ListRequiredStatusChecks_Call) Run(run func(ctx context.Context, _a1 id.PR, branch string)) *MockAPI_ListRequiredStatusChecks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR), args[2].(string))
	})
	return _c
}

func (_c *MockAPI_ListRequiredStatusChecks_Call) Return(_a0 []string, _a1 error) *MockAPI_ListRequiredStatusChecks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListRequiredStatusChecks_Call) RunAndReturn(run func(context.Context, id.PR, string) ([]string, error)) *MockAPI_ListRequiredStatusChecks_Call {
	_c.Call.Return(run)
	return _c
}

// ListReviews provides a mock function with given fields: ctx, _a1
func (_m *MockAPI) ListReviews(ctx context.Context, _a1 id.PR) ([]*v50github.PullRequestReview, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListReviews")
	}

	var r0 []*v50github.PullRequestReview
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) ([]*v50github.PullRequestReview, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.PR) []*v50github.PullRequestReview); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v50github.PullRequestReview)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.PR) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_ListReviews_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReviews'
type MockAPI_ListReviews_Call struct {
	*mock.Call
}

// ListReviews is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 id.PR
func (_e *MockAPI_Expecter) ListReviews(ctx interface{}, _a1 interface{}) *MockAPI_ListReviews_Call {
	return &MockAPI_ListReviews_Call{Call: _e.mock.On("ListReviews", ctx, _a1)}
}

func (_c *MockAPI_ListReviews_Call) Run(run func(ctx context.Context, _a1 id.PR)) *MockAPI_ListReviews_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(id.PR))
	})
	return _c
}

func (_c *MockAPI_ListReviews_Call) Return(_a0 []*v50github.PullRequestReview, _a1 error) *MockAPI_ListReviews_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_ListReviews_Call) RunAndReturn(run func(context.Context, id.PR) ([]*v50github.PullRequestReview, error)) *MockAPI_ListReviews_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAPI creates a new instance of MockAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAPI {
	mock := &MockAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
