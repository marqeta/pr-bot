package identity

import (
	"context"
	"errors"
)

type IdentityValidator interface {
	ValidateIdentity(ctx context.Context, identity *CallerIdentity) error
}

// Default validator: allows anything with a non-empty ARN
type AllowAllValidator struct{}

func (v *AllowAllValidator) ValidateIdentity(ctx context.Context, identity *CallerIdentity) error {
	if identity == nil || identity.Arn == "" {
		return errors.New("identity missing or invalid")
	}
	return nil
}
