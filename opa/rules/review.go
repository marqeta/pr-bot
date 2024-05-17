package rules

import (
	"context"
	"fmt"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
	"github.com/open-policy-agent/opa/sdk"
)

type Review struct {
	RuleName string
	client   client.Client
}

// Evaluate implements Rules.
func (r *Review) Evaluate(
	// TODO: use mapstructure pkg to convert arbitrary object into expected struct
	// parsing using mapstructure is slow,
	// therefore parsing individual fields since there is only two fields.
	ctx context.Context, module string, input *input.Model) (types.Review, error) {

	opt := sdk.DecisionOptions{
		Path:       r.client.Path(module, r.RuleName) + "/type",
		Input:      input,
		DecisionID: r.client.DecisionID(ctx),
	}

	typeDecision, err := r.client.Decision(ctx, opt)
	if err != nil {
		return types.Review{}, err
	}
	typeStr, ok := typeDecision.Result.(string)
	if !ok {
		return types.Review{}, fmt.Errorf("%w rule: review.type, expected: string, got: %T", ErrInvalidReturnType, typeDecision.Result)
	}
	reviewType, err := types.ParseReviewType(typeStr)
	if err != nil {
		return types.Review{}, err
	}

	if reviewType == types.Skip {
		return types.Review{
			Type: reviewType,
		}, nil
	}

	optBody := sdk.DecisionOptions{
		Path:       r.client.Path(module, r.RuleName) + "/body",
		Input:      input,
		DecisionID: r.client.DecisionID(ctx),
	}

	bodyDecision, err := r.client.Decision(ctx, optBody)
	if err != nil {
		return types.Review{}, err
	}

	body, ok := bodyDecision.Result.(string)
	if !ok {
		return types.Review{}, fmt.Errorf("%w rule: review.body, expected: string, got: %T", ErrInvalidReturnType, bodyDecision.Result)
	}

	optPref := sdk.DecisionOptions{
		Path:       r.client.Path(module, r.RuleName) + "/merge_preference",
		Input:      input,
		DecisionID: r.client.DecisionID(ctx),
	}
	prefDecision, err := r.client.Decision(ctx, optPref)
	if err != nil {
		return types.Review{}, err
	}
	prefStr, ok := prefDecision.Result.(string)
	if !ok {
		return types.Review{}, fmt.Errorf("%w rule: review.merge_preference, expected: string, got: %T", types.ErrInvalidMergePreference, prefDecision.Result)
	}
	pref, err := types.ParseMergeMethod(prefStr)
	if err != nil {
		return types.Review{}, err
	}

	return types.Review{
		Type:            reviewType,
		Body:            body,
		MergePreference: pref,
	}, nil
}

func NewReview(RuleName string, client client.Client) Rules[types.Review] {
	return &Review{
		RuleName: RuleName,
		client:   client,
	}
}
