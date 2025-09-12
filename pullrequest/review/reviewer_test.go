package review_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/shurcooL/githubv4"
)

func Test_reviewer_Approve(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	//nolint:goerr113
	errAutomerge := errors.New("pull request auto merge is not allowed")
	//nolint:goerr113
	errAutomergeUppercase := errors.New("Pull Request auto merge is not allowed")
	//nolint:goerr113
	errAutomergeHasHooks := errors.New("pull request is in has_hooks status")
	//nolint:goerr113
	errAutomergeHasHooksUppercase := errors.New("pull request is in has_hooks status")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(d *gh.MockAPI)
		mergeMethod     githubv4.PullRequestMergeMethod
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should auto merge PR",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(nil)
					d.EXPECT().AddReview(ctx, sampleID(), "LGTM", gh.Approve).
						Return(nil)
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: false,
		},
		{
			name: "Should auto merge PR with rebase",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodRebase).
						Return(nil)
					d.EXPECT().AddReview(ctx, sampleID(), "LGTM", gh.Approve).
						Return(nil)
				},
				mergeMethod: githubv4.PullRequestMergeMethodRebase,
			},
			wantErr: false,
		},
		{
			name: "Should auto merge PR with squash",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodSquash).
						Return(nil)
					d.EXPECT().AddReview(ctx, sampleID(), "LGTM", gh.Approve).
						Return(nil)
				},
				mergeMethod: githubv4.PullRequestMergeMethodSquash,
			},
			wantErr: false,
		},
		{
			name: "Throw error when EnableAutoMerge fails",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(errRandom).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
		{
			name: "Throw error when AutoMerge is not allowed in repo",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(errAutomerge).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
		{
			name: "Throw error when AutoMerge is not allowed in repo, case sensitive error",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(errAutomergeUppercase).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
		{
			name: "Throw error when pr is in has_hooks status",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(errAutomergeHasHooks).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
		{
			name: "Throw error when pr is in has_hooks status uppercase",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(errAutomergeHasHooksUppercase).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
		{
			name: "Throw error when AddReview fails",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).
						Return(nil).Once()
					d.EXPECT().AddReview(ctx, sampleID(), "LGTM", gh.Approve).
						Return(errRandom).Once()
				},
				mergeMethod: githubv4.PullRequestMergeMethodMerge,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := gh.NewMockAPI(t)
			metrics := metrics.NewNoopEmitter()
			r := review.NewReviewer(mockAPI, metrics)
			tt.args.setExpectations(mockAPI)
			if err := r.Approve(ctx, tt.args.id, tt.args.body, review.ApproveOptions{
				MergeMethod: tt.args.mergeMethod,
			}); (err != nil) != tt.wantErr {
				t.Errorf("reviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_reviewer_Approve_Sleep(t *testing.T) {
	ctx := context.Background()
	mockAPI := gh.NewMockAPI(t)
	metrics := metrics.NewNoopEmitter()
	r := review.NewReviewer(mockAPI, metrics)

	// Set up expectations
	mockAPI.EXPECT().EnableAutoMerge(ctx, sampleID(), githubv4.PullRequestMergeMethodMerge).Return(nil)
	mockAPI.EXPECT().AddReview(ctx, sampleID(), "LGTM", gh.Approve).Return(nil)

	// Measure the time to verify sleep
	start := time.Now()
	err := r.Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
		MergeMethod: githubv4.PullRequestMergeMethodMerge,
	})
	duration := time.Since(start)

	// Verify no error and sleep duration is at least 10 second
	if err != nil {
		t.Errorf("Approve() error = %v, want nil", err)
	}
	if duration < 10*time.Second {
		t.Errorf("Sleep duration = %v, want at least 10s", duration)
	}
}

func Test_reviewer_Comment(t *testing.T) {
	ctx := context.Background()

	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(d *gh.MockAPI)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should comment on a PR",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().AddReview(ctx, sampleID(), "random body", gh.Comment).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Throw error when AddReview fails",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().AddReview(ctx, sampleID(), "random body", gh.Comment).
						Return(errRandom).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := gh.NewMockAPI(t)
			metrics := metrics.NewNoopEmitter()
			r := review.NewReviewer(mockAPI, metrics)
			tt.args.setExpectations(mockAPI)
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("reviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_reviewer_RequestChanges(t *testing.T) {
	ctx := context.Background()

	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(d *gh.MockAPI)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should request changes on PR",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().AddReview(ctx, sampleID(), "random body", gh.RequestChanges).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Throw error when AddReview fails",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().AddReview(ctx, sampleID(), "random body", gh.RequestChanges).
						Return(errRandom).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := gh.NewMockAPI(t)
			metrics := metrics.NewNoopEmitter()
			r := review.NewReviewer(mockAPI, metrics)
			tt.args.setExpectations(mockAPI)
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("reviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func sampleID() id.PR {

	return id.PR{
		Owner:        "owner1",
		Repo:         "repo1",
		Number:       1,
		NodeID:       "nodeid1",
		RepoFullName: "owner1/repo1",
		Author:       "user1",
		URL:          "owner1/repo1/1",
	}
}
