package identity

import (
	"errors"
)

type Validator interface {
	ValidateIdentity(identity *CallerIdentity) error
}

// Default validator: allows anything with a non-empty ARN
type AllowAllValidator struct{}

func (v *AllowAllValidator) ValidateIdentity(identity *CallerIdentity) error {
	if identity == nil || identity.Arn == "" {
		//nolint:goerr113
		return errors.New("identity missing or invalid")
	}
	return nil
}
