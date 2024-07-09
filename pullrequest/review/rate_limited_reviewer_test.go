package review_test

import (
	"context"
	"errors"
	"testing"

	pe "github.com/marqeta/pr-bot/errors"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/marqeta/pr-bot/rate"
	lim "github.com/mennanov/limiters"
	"github.com/shurcooL/githubv4"
)

func Test_rateLimitedReviewer_Approve(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	errThrottled := pe.TooManyRequestError(ctx, "throttled", lim.ErrLimitExhausted)
	type args struct {
		id              id.PR
		body            string
		throttlerErr    error
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate.Approve when throttler does not return error",
			args: args{
				id:           sampleID(),
				body:         "random body",
				throttlerErr: nil,
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().Approve(ctx, sampleID(), "random body", review.ApproveOptions{
						MergeMethod: githubv4.PullRequestMergeMethodMerge,
					}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error when throttler returns error",
			args: args{
				id:           sampleID(),
				body:         "random body",
				throttlerErr: errRandom,
				setExpectations: func(_ *gh.MockAPI, _ *review.MockReviewer) {
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when throttler returns error",
			args: args{
				id:           sampleID(),
				body:         "random body",
				throttlerErr: errThrottled,
				setExpectations: func(_ *gh.MockAPI, _ *review.MockReviewer) {
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			throttler := rate.NewMockThrottler(tt.args.throttlerErr)
			tt.args.setExpectations(api, delegate)
			r := review.NewRateLimitedReviewer(delegate, api, throttler)
			if err := r.Approve(ctx, tt.args.id, tt.args.body, review.ApproveOptions{
				MergeMethod: githubv4.PullRequestMergeMethodMerge,
			}); (err != nil) != tt.wantErr {
				t.Errorf("rateLimitedReviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rateLimitedReviewer_Comment(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate.Comment",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().Comment(ctx, sampleID(), "random body").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error from delegate.Comment",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().Comment(ctx, sampleID(), "random body").
						Return(errRandom)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			throttler := rate.NewMockThrottler(nil)
			tt.args.setExpectations(api, delegate)
			r := review.NewRateLimitedReviewer(delegate, api, throttler)
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("rateLimitedReviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rateLimitedReviewer_RequestChanges(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate.RequestChanges",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "random body").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error from delegate.RequestChanges",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "random body").
						Return(errRandom)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			throttler := rate.NewMockThrottler(nil)
			tt.args.setExpectations(api, delegate)
			r := review.NewRateLimitedReviewer(delegate, api, throttler)
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("rateLimitedReviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
