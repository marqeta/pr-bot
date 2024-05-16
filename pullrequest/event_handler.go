package pullrequest

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/google/go-github/v50/github"
	"github.com/shurcooL/githubv4"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
	"github.com/marqeta/pr-bot/pullrequest/review"
)

//go:generate mockery --name EventHandler --testonly
type EventHandler interface {
	EvalAndReview(ctx context.Context, id id.PR, event *github.PullRequestEvent) error
}

type eventHandler struct {
	reviewer  review.Reviewer
	metrics   metrics.Emitter
	evaluator opa.Evaluator
}

func NewEventHandler(evaluator opa.Evaluator, reviewer review.Reviewer, metrics metrics.Emitter) EventHandler {
	return &eventHandler{
		reviewer:  reviewer,
		metrics:   metrics,
		evaluator: evaluator,
	}
}

func (eh *eventHandler) EvalAndReview(ctx context.Context, id id.PR, event *github.PullRequestEvent) error {
	oplog := httplog.LogEntry(ctx)

	ghe := input.ToGHE(event)
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
		return eh.reviewer.Approve(ctx, id, opaResult.Review.Body, review.ApproveOptions{
			AutoMergeEnabled: event.GetPullRequest().GetBase().GetRepo().GetAllowAutoMerge(),
			DefaultBranch:    event.Repo.GetDefaultBranch(),
			MergeMethod:      mergeMethod(event),
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

func mergeMethod(event *github.PullRequestEvent) githubv4.PullRequestMergeMethod {
	rebase := event.PullRequest.GetBase().GetRepo().GetAllowRebaseMerge()
	squash := event.PullRequest.GetBase().GetRepo().GetAllowSquashMerge()

	if rebase {
		return githubv4.PullRequestMergeMethodRebase
	}
	if squash {
		return githubv4.PullRequestMergeMethodSquash
	}
	return githubv4.PullRequestMergeMethodMerge

}
