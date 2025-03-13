package plugins

import (
	"context"
	"encoding/json"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
)

type PullRequestReviewers struct {
	dao gh.API
}

// GetInputMsg implements input.Plugin.
func (prr *PullRequestReviewers) GetInputMsg(ctx context.Context, ghe input.GHE) (json.RawMessage, error) {
	reviews, err := prr.dao.ListReviews(ctx, ghe.ToID())
	if err != nil {
		return json.RawMessage{}, err
	}
	data, err := json.Marshal(reviews)
	if err != nil {
		return json.RawMessage{}, err
	}
	return json.RawMessage(data), nil
}

// Name implements input.Plugin.
func (prr *PullRequestReviewers) Name() string {
	return "reviews"
}

func NewPullRequestReviewers(dao gh.API) input.Plugin {
	return &PullRequestReviewers{dao: dao}
}
