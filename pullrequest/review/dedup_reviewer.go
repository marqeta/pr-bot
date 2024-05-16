package review

import (
	"context"

	"github.com/go-chi/httplog"
	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/opa/types"
)

const (
	ListReviewsError = "error listing reviews on PR %v"
)

type DedupReviewer struct {
	api            gh.API
	delegate       Reviewer
	serviceAccount string
}

// Approve implements Reviewer.
func (d *DedupReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)
	reviews, err := d.api.ListReviews(ctx, id)
	if err != nil {
		oplog.Err(err).Msgf(ListReviewsError, id.URL)
		return err
	}
	if d.checkForReview(reviews, types.Approve) {
		oplog.Info().Msgf("PR already has a review of type %v or higher", types.Approve)
		return nil
	}
	return d.delegate.Approve(ctx, id, body, opts)
}

// Comment implements Reviewer.
func (d *DedupReviewer) Comment(ctx context.Context, id id.PR, body string) error {
	oplog := httplog.LogEntry(ctx)
	reviews, err := d.api.ListReviews(ctx, id)
	if err != nil {
		oplog.Err(err).Msgf(ListReviewsError, id.URL)
		return err
	}
	if d.checkForReview(reviews, types.Comment) {
		oplog.Info().Msgf("PR already has a review of type %v or higher", types.Comment)
		return nil
	}
	return d.delegate.Comment(ctx, id, body)
}

// RequestChanges implements Reviewer.
func (d *DedupReviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	oplog := httplog.LogEntry(ctx)
	reviews, err := d.api.ListReviews(ctx, id)
	if err != nil {
		oplog.Err(err).Msgf(ListReviewsError, id.URL)
		return err
	}
	if d.checkForReview(reviews, types.RequestChanges) {
		oplog.Info().Msgf("PR already has a review of type %v or higher", types.RequestChanges)
		return nil
	}
	return d.delegate.RequestChanges(ctx, id, body)
}

// NewDedupReviewer returns a new DedupReviewer
func NewDedupReviewer(delegate Reviewer, api gh.API, serviceAccount string) Reviewer {
	return &DedupReviewer{
		api:            api,
		delegate:       delegate,
		serviceAccount: serviceAccount,
	}
}

// checkForReview checks if the service account has already reviewed the PR with a review type of minReviewType or higher
func (d *DedupReviewer) checkForReview(reviews []*github.PullRequestReview, minReviewType types.ReviewType) bool {
	for _, review := range reviews {
		if review.User.GetLogin() == d.serviceAccount {
			t, err := types.ParseReviewState(review.GetState())
			if err != nil {
				continue
			}
			if t >= minReviewType {
				return true
			}
		}
	}
	return false
}
