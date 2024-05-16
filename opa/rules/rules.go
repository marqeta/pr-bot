package rules

import (
	"context"
	"errors"

	"github.com/marqeta/pr-bot/opa/input"
)

var ErrInvalidReturnType = errors.New("invalid return type from rule evaluation")

//go:generate mockery --name Rules
type Rules[T any] interface {
	Evaluate(ctx context.Context, module string, input *input.Model) (T, error)
}
