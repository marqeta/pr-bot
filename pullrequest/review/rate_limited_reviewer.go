package review

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/httplog"
	pe "github.com/marqeta/pr-bot/errors"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/rate"
)

type rateLimitedReviewer struct {
	api       gh.API
	delegate  Reviewer
	throttler rate.Throttler
}

// Approve implements Reviewer.
func (r *rateLimitedReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)
	err := r.throttler.ShouldThrottle(ctx, id)
	if err != nil {
		var ae pe.APIError
		if errors.As(err, &ae) && ae.StatusCode == http.StatusTooManyRequests {
			oplog.Err(ae).Msgf("request throttled for PR %v", id.URL)
			// TODO publish error in UI and/or as comments on PR
		} else {
			oplog.Err(err).Msgf("error from throttler for PR %v", id.URL)
		}
		return err
	}
	return r.delegate.Approve(ctx, id, body, opts)
}

// Comment implements Reviewer.
func (r *rateLimitedReviewer) Comment(ctx context.Context, id id.PR, body string) error {
	return r.delegate.Comment(ctx, id, body)
}

// RequestChanges implements Reviewer.
func (r *rateLimitedReviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	return r.delegate.RequestChanges(ctx, id, body)
}

func NewRateLimitedReviewer(delegate Reviewer, api gh.API, throttler rate.Throttler) Reviewer {
	return &rateLimitedReviewer{
		delegate:  delegate,
		api:       api,
		throttler: throttler,
	}
}

func (r *rateLimitedReviewer) Dismiss(ctx context.Context, id id.PR, body string) error {
	return r.delegate.Dismiss(ctx, id, body)
}
