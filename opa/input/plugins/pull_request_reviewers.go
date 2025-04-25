package plugins

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v50/github"

	//"github.com/google/go-github/v50/github"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
)

type PullRequestReviewers struct {
	dao gh.API
}

func (prr *PullRequestReviewers) GetInputMsg(ctx context.Context, ghe input.GHE) (json.RawMessage, error) {
	// 1. Fetch all reviews (note: []*github.PullRequestReview)
	reviews, err := prr.dao.ListReviews(ctx, ghe.ToID())
	if err != nil {
		return nil, err
	}

	// 2. Reverse-iterate and pick only the first occurrence per login
	//    key = string login, value = most recent *github.PullRequestReview
	latestByUser := make(map[string]*github.PullRequestReview)
	for i := len(reviews) - 1; i >= 0; i-- {
		r := reviews[i]
		login := *r.User.Login
		if _, seen := latestByUser[login]; !seen {
			latestByUser[login] = r
		}
	}

	// 3. Turn the map back into a slice
	latestReviews := make([]*github.PullRequestReview, 0, len(latestByUser))
	for _, r := range latestByUser {
		latestReviews = append(latestReviews, r)
	}

	// 4. Marshal exactly that slice of “latest” reviews
	data, err := json.Marshal(latestReviews)
	if err != nil {
		return nil, err
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
