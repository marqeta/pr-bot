package client

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/open-policy-agent/opa/v1/sdk"
)

//go:generate mockery --name Client
type Client interface {
	Decision(ctx context.Context, options sdk.DecisionOptions) (*sdk.DecisionResult, error)
	DecisionID(ctx context.Context) string
	Path(module string, rule string) string
}

type client struct {
	sdk *sdk.OPA
}

func (c *client) DecisionID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func (c *client) Path(module string, rule string) string {
	return fmt.Sprintf("%s/%s", module, rule)
}

// Decision implements Client.
func (c *client) Decision(ctx context.Context, options sdk.DecisionOptions) (*sdk.DecisionResult, error) {
	return c.sdk.Decision(ctx, options)
}

func NewClient(sdk *sdk.OPA) Client {
	return &client{
		sdk: sdk,
	}
}
