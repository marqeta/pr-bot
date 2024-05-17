package rules

import (
	"context"
	"fmt"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/open-policy-agent/opa/sdk"
)

type Track struct {
	RuleName string
	client   client.Client
}

// Evaluate implements Rules.
func (t *Track) Evaluate(
	ctx context.Context, module string, input *input.Model) (bool, error) {

	opt := sdk.DecisionOptions{
		Path:       t.client.Path(module, t.RuleName),
		Input:      input,
		DecisionID: t.client.DecisionID(ctx),
	}

	result, err := t.client.Decision(ctx, opt)
	if err != nil {
		return false, err
	}

	track, ok := result.Result.(bool)
	if !ok {
		return false, fmt.Errorf("%w rule: track, expected: bool, got: %T", ErrInvalidReturnType, result.Result)
	}
	return track, nil
}

func NewTrack(ruleName string, client client.Client) Rules[bool] {
	return &Track{
		RuleName: ruleName,
		client:   client,
	}
}
