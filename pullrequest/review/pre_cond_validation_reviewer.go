package review

import (
	"context"
	"errors"

	"github.com/go-chi/httplog"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/id"
)

var ErrAutoMergeDisabled = errors.New("auto merge is disabled in repo")

type preCondValidationReviewer struct {
	delegate Reviewer
}

// Approve implements Reviewer.
func (p *preCondValidationReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)

	if !opts.AutoMergeEnabled {
		ae := pe.UserError(ctx, AutoMergeError, ErrAutoMergeDisabled)
		oplog.Error().Msgf("Auto merge is disabled in repo for pr %v", id.URL)
		return ae
	}

	return p.delegate.Approve(ctx, id, body, opts)
}

// Comment implements Reviewer.
func (p *preCondValidationReviewer) Comment(ctx context.Context, id id.PR, body string) error {
	return p.delegate.Comment(ctx, id, body)
}

// RequestChanges implements Reviewer.
func (p *preCondValidationReviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	return p.delegate.RequestChanges(ctx, id, body)
}

func NewPreCondValidationReviewer(delegate Reviewer) Reviewer {
	return &preCondValidationReviewer{
		delegate: delegate,
	}
}
