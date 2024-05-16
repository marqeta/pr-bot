package input

import (
	"context"
	"encoding/json"
)

//go:generate mockery --name Plugin --testonly
type Plugin interface {
	GetInputMsg(ctx context.Context, ghe GHE) (json.RawMessage, error)
	Name() string
}
