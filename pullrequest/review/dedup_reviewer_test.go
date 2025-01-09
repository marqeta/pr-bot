package review_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/pullrequest/review"
)

func TestDedupReviewer_Approve(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		serviceAccount  string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should approve PR for the first time",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PR for the first time empty list of reviews",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PR even with reviews from other users",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "approved")}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not re-approve PR",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("approved"), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PR when there is a comment",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("commented"), nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PR when there is a request for changes",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("changes_requested"), nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not approve PR when list reviews returns an error",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should approve PR when review type is not parsable",
			args: args{
				id:             sampleID(),
				body:           "LGTM",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("asd"), nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{}).
						Return(nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewDedupReviewer(delegate, api, tt.args.serviceAccount)
			if err := r.Approve(ctx, tt.args.id, tt.args.body, review.ApproveOptions{}); (err != nil) != tt.wantErr {
				t.Errorf("DedupReviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDedupReviewer_Comment(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		serviceAccount  string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should comment on PR for the first time",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PR for the first time empty list of reviews",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{}, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PR even with approval from other users",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "approved")}, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PR even with request_changes from other users",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "changes_requested")}, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not re-comment on PR",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("commented"), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PR when there is approval",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("approved"), nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PR when there is a request for changes",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("changes_requested"), nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not comment PR when list reviews returns an error",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should comment PR when review type is not parsable",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("asd"), nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").
						Return(nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewDedupReviewer(delegate, api, tt.args.serviceAccount)
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("DedupReviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDedupReviewer_RequestChanges(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		serviceAccount  string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should request_changes on PR for the first time",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR for the first time empty list of reviews",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{}, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").
						Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR even with approval from other users",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "approved")}, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR even with comments from other users",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "commented")}, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR even with request_changes from other users",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return([]*github.PullRequestReview{prReview("asd", "changes_requested")}, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not re-request_changes on PR",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("changes_requested"), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR when there is a comment",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("commented"), nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request_changes on PR when there is a approval",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("approved"), nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not request_changes on PR when list reviews returns an error",
			args: args{
				id:             sampleID(),
				body:           "comment",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should request_changes PR when review type is not parsable",
			args: args{
				id:             sampleID(),
				body:           "request_changes",
				serviceAccount: "svc-ci-prbot",
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListReviews(ctx, sampleID()).
						Return(reviews("asd"), nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "request_changes").
						Return(nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewDedupReviewer(delegate, api, tt.args.serviceAccount)
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("DedupReviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func reviews(states ...string) []*github.PullRequestReview {
	var r []*github.PullRequestReview
	for _, state := range states {
		r = append(r, prReview("svc-ci-prbot", state))
	}
	return r
}

func prReview(user, state string) *github.PullRequestReview {
	return &github.PullRequestReview{
		User: &github.User{
			Login: github.String(user),
		},
		State: github.String(state),
	}
}
