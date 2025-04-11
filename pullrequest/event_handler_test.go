package pullrequest_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
	"github.com/marqeta/pr-bot/pullrequest"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/shurcooL/githubv4"
)

func Test_eventHandlerV2_EvalAndReview(t *testing.T) {
	ctx := context.TODO()
	type args struct {
		id             id.PR
		event          *github.PullRequestEvent
		setExpectaions func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should approve PRs",
			args: args{
				id:    sampleID(),
				event: merge(prEvent(github.String("opened"), sampleID())),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: true,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodMerge,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PRs with auto merge disabled",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: false,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodMerge,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PRs with random default branch",
			args: args{
				id:    sampleID(),
				event: randomBranchAsDefault(prEvent(github.String("opened"), sampleID())),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: false,
						DefaultBranch:    "random",
						MergeMethod:      githubv4.PullRequestMergeMethodMerge,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should call approve with rebase merge method",
			args: args{
				id:    sampleID(),
				event: rebase(prEvent(github.String("opened"), sampleID())),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: true,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodRebase,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PRs with squash merge method for empty commits",
			args: args{
				id:    sampleID(),
				event: empty(allMergeMethods(prEvent(github.String("opened"), sampleID()))),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: true,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodSquash,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PRs with merge commit for empty commits",
			args: args{
				id:    sampleID(),
				event: empty(rebase(merge(prEvent(github.String("opened"), sampleID())))),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: true,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodMerge,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should approve PRs with squash merge method",
			args: args{
				id:    sampleID(),
				event: squash(prEvent(github.String("opened"), sampleID())),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: true,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodSquash,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should request changes on PRs",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.RequestChanges,
							Body: "Change me",
						},
					}, nil)
					r.EXPECT().RequestChanges(ctx, sampleID(), "Change me").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should comment on PRs",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Comment,
							Body: "Comment",
						},
					}, nil)
					r.EXPECT().Comment(ctx, sampleID(), "Comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not review PRs when track=false",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, _ *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: false,
					}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error when when evaluation fails",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, _ *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{}, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error when when evaluation fails",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, _ *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{}, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error when approval fails",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Approve,
							Body: "LGTM",
						},
					}, nil)
					r.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						AutoMergeEnabled: false,
						DefaultBranch:    "main",
						MergeMethod:      githubv4.PullRequestMergeMethodMerge,
					}).Return(errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error when request changes fails",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.RequestChanges,
							Body: "Change me",
						},
					}, nil)
					r.EXPECT().RequestChanges(ctx, sampleID(), "Change me").Return(errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error when comment fails",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, r *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Comment,
							Body: "Comment",
						},
					}, nil)
					r.EXPECT().Comment(ctx, sampleID(), "Comment").Return(errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should skip PRs with no action when result.type is skip",
			args: args{
				id:    sampleID(),
				event: prEvent(github.String("opened"), sampleID()),
				setExpectaions: func(e *opa.MockEvaluator, _ *review.MockReviewer, p *github.PullRequestEvent) {
					e.EXPECT().Evaluate(ctx, ToGHE(p)).Return(types.Result{
						Track: true,
						Review: types.Review{
							Type: types.Skip,
						},
					}, nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := opa.NewMockEvaluator(t)
			r := review.NewMockReviewer(t)
			m := metrics.NewNoopEmitter()
			a := input.NewMockAdapter(t)
			eh := pullrequest.NewEventHandler(e, r, m, a)
			tt.args.setExpectaions(e, r, tt.args.event)
			ghe := ToGHE(tt.args.event)
			if err := eh.EvalAndReview(ctx, tt.args.id, ghe); (err != nil) != tt.wantErr {
				t.Errorf("eventHandlerV2.Review() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ToGHE(event *github.PullRequestEvent) input.GHE {
	return input.GHE{
		Event:        "pull_request",
		Action:       event.GetAction(),
		PullRequest:  event.GetPullRequest(),
		Repository:   event.GetRepo(),
		Organization: event.GetOrganization(),
	}
}
