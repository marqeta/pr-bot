package pullrequest

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/datastore"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/shurcooL/githubv4"
)

//go:generate mockery --name EventHandler
type EventHandler interface {
	EvalAndReview(ctx context.Context, id id.PR, ghe input.GHE) error
	EvalAndReviewPREvent(ctx context.Context, id id.PR, event *github.PullRequestEvent) error
	EvalAndReviewPRReviewEvent(ctx context.Context, id id.PR, event *github.PullRequestReviewEvent) error
	EvalAndReviewDataEvent(ctx context.Context, metadata *datastore.Metadata) error
}

type eventHandler struct {
	reviewer  review.Reviewer
	metrics   metrics.Emitter
	evaluator opa.Evaluator
	adapter   input.Adapter
}

func NewEventHandler(evaluator opa.Evaluator, reviewer review.Reviewer, metrics metrics.Emitter, adapter input.Adapter) EventHandler {
	return &eventHandler{
		reviewer:  reviewer,
		metrics:   metrics,
		evaluator: evaluator,
		adapter:   adapter,
	}
}

func (eh *eventHandler) EvalAndReviewDataEvent(ctx context.Context, metadata *datastore.Metadata) error {
	ghe, err := eh.adapter.MetadataToGHE(ctx, metadata)
	if err != nil {
		return err
	}
	return eh.EvalAndReview(ctx, ghe.ToID(), ghe)
}

func (eh *eventHandler) EvalAndReviewPREvent(ctx context.Context, id id.PR, event *github.PullRequestEvent) error {
	ghe, err := eh.adapter.PREventToGHE(ctx, event)
	if err != nil {
		return err
	}
	return eh.EvalAndReview(ctx, id, ghe)
}

func (eh *eventHandler) EvalAndReviewPRReviewEvent(ctx context.Context, id id.PR, event *github.PullRequestReviewEvent) error {
	ghe, err := eh.adapter.PRReviewEventToGHE(ctx, event)
	if err != nil {
		return err
	}
	return eh.EvalAndReview(ctx, id, ghe)
}

func (eh *eventHandler) EvalAndReview(ctx context.Context, id id.PR, ghe input.GHE) error {
	oplog := httplog.LogEntry(ctx)

	tags := id.ToTags()
	opaResult, err := eh.evaluator.Evaluate(ctx, ghe)
	oplog.Err(err).Interface("decision", opaResult).Msg("opa evaluation complete")
	if err != nil {
		eh.metrics.EmitDist(ctx, "opa.evaluator.errors", 1.0, tags)
		return err
	}

	if !opaResult.Track {
		oplog.Info().Msg("track=false, skipping review")
		return nil
	}

	switch opaResult.Review.Type {

	case types.Approve:
		autoApprove := ghe.PullRequest.GetBase().GetRepo().GetAllowAutoMerge() ||
			ghe.Repository.GetAllowAutoMerge()

		return eh.reviewer.Approve(ctx, id, opaResult.Review.Body, review.ApproveOptions{
			AutoMergeEnabled: autoApprove,
			DefaultBranch:    ghe.Repository.GetDefaultBranch(),
			MergeMethod:      mergeMethod(ghe),
		})
	case types.RequestChanges:
		return eh.reviewer.RequestChanges(ctx, id, opaResult.Review.Body)
	case types.Comment:
		return eh.reviewer.Comment(ctx, id, opaResult.Review.Body)
	default:
		oplog.Info().Msg("skipping review")
	}
	return nil
}

func mergeMethod(ghe input.GHE) githubv4.PullRequestMergeMethod {
	rebase := ghe.PullRequest.GetBase().GetRepo().GetAllowRebaseMerge() || ghe.Repository.GetAllowRebaseMerge()
	squash := ghe.PullRequest.GetBase().GetRepo().GetAllowSquashMerge() || ghe.Repository.GetAllowSquashMerge()
	fc := ghe.PullRequest.GetChangedFiles()
	// TODO: let policy specify what merge method to use.
	// when rebasing empty commits on to main,
	// no new commit is created, therefore no triggers would be fired.
	// use squash to force a new commit to be created. when merging empty PRs
	if rebase && fc > 0 {
		return githubv4.PullRequestMergeMethodRebase
	}
	if squash {
		return githubv4.PullRequestMergeMethodSquash
	}
	return githubv4.PullRequestMergeMethodMerge

}
