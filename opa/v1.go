package opa

import (
	"context"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/marqeta/pr-bot/opa/types"
)

const (
	TrackRuleName  = "track"
	ReviewRuleName = "review"
)

type V1 struct {
	track  rules.Rules[bool]
	review rules.Rules[types.Review]
}

func NewV1PolicyFromRules(t rules.Rules[bool], r rules.Rules[types.Review]) Policy {
	return &V1{
		track:  t,
		review: r,
	}
}

func NewV1Policy(client client.Client) Policy {
	return NewV1PolicyFromRules(
		rules.NewTrack(TrackRuleName, client),
		rules.NewReview(ReviewRuleName, client),
	)
}

// Evaluate implements Policy.
func (v1 *V1) Evaluate(
	ctx context.Context, module string, input *input.Model) (types.Result, error) {

	track, err := v1.track.Evaluate(ctx, module, input)
	if err != nil {
		return types.Result{}, err
	}
	if !track {
		return types.Result{
			Track: track,
		}, nil
	}

	review, err := v1.review.Evaluate(ctx, module, input)
	if err != nil {
		return types.Result{}, err
	}
	return types.Result{
		Track:  track,
		Review: review,
	}, nil
}
