package input

import (
	"context"

	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/datastore"
	gh "github.com/marqeta/pr-bot/github"
)

//go:generate mockery --name Adapter
type Adapter interface {
	PREventToGHE(ctx context.Context, event *github.PullRequestEvent) (GHE, error)
	MetadataToGHE(ctx context.Context, metadata *datastore.Metadata) (GHE, error)
	PRReviewEventToGHE(ctx context.Context, event *github.PullRequestReviewEvent) (GHE, error)
}

type adapter struct {
	dao gh.API
}

func NewAdapter(dao gh.API) Adapter {
	return &adapter{dao: dao}
}

func (a *adapter) PREventToGHE(_ context.Context, event *github.PullRequestEvent) (GHE, error) {
	return GHE{
		Event:        "pull_request",
		Action:       event.GetAction(),
		PullRequest:  event.GetPullRequest(),
		Repository:   event.GetRepo(),
		Organization: event.GetOrganization(),
	}, nil
}

func (a *adapter) MetadataToGHE(ctx context.Context, metadata *datastore.Metadata) (GHE, error) {
	org, err := a.dao.GetOrganization(ctx, metadata.PR)
	if err != nil {
		return GHE{}, err
	}
	repo, err := a.dao.GetRepository(ctx, metadata.PR)
	if err != nil {
		return GHE{}, err
	}
	pr, err := a.dao.GetPullRequest(ctx, metadata.PR)
	if err != nil {
		return GHE{}, err
	}

	return GHE{
		Event:        metadata.Service,
		Action:       metadata.Job,
		PullRequest:  pr,
		Repository:   repo,
		Organization: org,
	}, nil
}

func (a *adapter) PRReviewEventToGHE(_ context.Context, event *github.PullRequestReviewEvent) (GHE, error) {
	return GHE{
		Event:        "pull_request_review",
		Action:       event.GetAction(),
		PullRequest:  event.GetPullRequest(),
		Repository:   event.GetRepo(),
		Organization: event.GetOrganization(),
	}, nil
}
