package review

import (
	"context"
	"errors"

	"github.com/go-chi/httplog"
	prbot "github.com/marqeta/pr-bot"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
)

var ErrAutoMergeDisabled = errors.New("auto merge is disabled in repo")

type preCondValidationReviewer struct {
	api      gh.API
	delegate Reviewer
	metrics  metrics.Emitter
}

// Approve implements Reviewer.
func (p *preCondValidationReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)

	if !opts.AutoMergeEnabled {
		ae := prbot.UserError(ctx, AutoMergeError, ErrAutoMergeDisabled)
		oplog.Error().Msgf("Auto merge is disabled in repo for pr %v", id.URL)
		return ae
	}

	checks, err := p.api.ListRequiredStatusChecks(ctx, id, opts.DefaultBranch)
	if err != nil {
		oplog.Err(err).Msgf("error listing required status checks on repo for PR %v", id.URL)
		ae := prbot.ServiceFault(ctx, "Error listing status checks on Repo", err)
		return ae
	}

	for _, check := range checks {
		if check == "blackbird-ci" {
			// blackbird-ci check is required, automerge and approve
			return p.delegate.Approve(ctx, id, body, opts)
		}
	}

	// No required blackbird-ci check, check if it is a BB build repo
	files, err := p.api.ListFilesInRootDir(ctx, id, opts.DefaultBranch)
	if err != nil {
		oplog.Err(err).Msgf("error listing files on repo for PR %v", id.URL)
		ae := prbot.ServiceFault(ctx, "Error listing files on Repo", err)
		return ae
	}

	for _, file := range files {
		if file == "blackbird.yaml" {
			// BB build repo, skip auto merge to avoid merging PR which doesnt pass blackbird-ci check
			oplog.Info().Msgf("blackbird-ci status check is not required on a blackbird repo skipping %v", id.URL)
			p.metrics.EmitDist(ctx, "noRequiredBBCheck", 1.0, id.ToTags())
			return nil
		}
	}
	// Not a BB build repo, auto merge and approve
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

func NewPreCondValidationReviewer(delegate Reviewer, api gh.API, metrics metrics.Emitter) Reviewer {
	return &preCondValidationReviewer{
		delegate: delegate,
		api:      api,
		metrics:  metrics,
	}
}
