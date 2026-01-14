package rules

import (
	"context"
	"fmt"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/open-policy-agent/opa/v1/sdk"
)

type Schema struct {
	RuleName string
	client   client.Client
}

// Evaluate implements Rules.
func (s *Schema) Evaluate(
	ctx context.Context, module string, input *input.Model) (string, error) {

	opt := sdk.DecisionOptions{
		Path:       s.client.Path(module, s.RuleName),
		Input:      input,
		DecisionID: s.client.DecisionID(ctx),
	}

	result, err := s.client.Decision(ctx, opt)
	if err != nil {
		return "", err
	}
	schema, ok := result.Result.(string)
	if !ok {
		return "", fmt.Errorf("%w rule: schema, expected: string, got: %T", ErrInvalidReturnType, result.Result)
	}
	return schema, nil
}

func NewSchema(name string, client client.Client) Rules[string] {
	return &Schema{
		RuleName: name,
		client:   client,
	}
}
