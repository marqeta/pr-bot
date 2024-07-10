package review

import (
	"context"
	"errors"

	"github.com/go-chi/httplog"
	pe "github.com/marqeta/pr-bot/errors"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
)

type bbReviewer struct {
	api      gh.API
	delegate Reviewer
	metrics  metrics.Emitter
}

var ErrBlackbirdCheckNotRequired = errors.New("blackbird-ci check is not a required status check on a blackbird repo")

// Approve implements Reviewer.
func (bb *bbReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)

	isBlackbirdRepo := false
	files, err := bb.api.ListFilesInRootDir(ctx, id, opts.DefaultBranch)
	if err != nil {
		oplog.Err(err).Send()
		return err
	}

	for _, file := range files {
		if file == "blackbird.yaml" {
			isBlackbirdRepo = true
			break
		}
	}

	if !isBlackbirdRepo {
		// not a blackbird repo, approve and automerge
		return bb.delegate.Approve(ctx, id, body, opts)
	}

	checks, err := bb.api.ListRequiredStatusChecks(ctx, id, opts.DefaultBranch)
	if err != nil {
		oplog.Err(err).Send()
		return err
	}

	for _, check := range checks {
		if check == "blackbird-ci" {
			// blackbird-ci check is required, automerge and approve
			return bb.delegate.Approve(ctx, id, body, opts)
		}
	}

	// BB built repo, skip auto merge to avoid merging PR which doesnt pass blackbird-ci check
	oplog.Info().Msgf("blackbird-ci status check is not required on a blackbird repo skipping %v", id.URL)
	bb.metrics.EmitDist(ctx, "noRequiredBBCheck", 1.0, id.ToTags())
	return pe.UserError(ctx, "blackbird-ci check must be required check", ErrBlackbirdCheckNotRequired)
}

// Comment implements Reviewer.
func (bb *bbReviewer) Comment(ctx context.Context, id id.PR, body string) error {
	return bb.delegate.Comment(ctx, id, body)
}

// RequestChanges implements Reviewer.
func (bb *bbReviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	return bb.delegate.RequestChanges(ctx, id, body)
}

func NewBBReviewer(delegate Reviewer, api gh.API, metrics metrics.Emitter) Reviewer {
	return &bbReviewer{
		delegate: delegate,
		api:      api,
		metrics:  metrics,
	}
}
