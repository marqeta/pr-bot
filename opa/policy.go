package opa

import (
	"context"

	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
)

//go:generate mockery --name Policy --testonly
type Policy interface {
	Evaluate(ctx context.Context, module string, input *input.Model) (types.Result, error)
}
