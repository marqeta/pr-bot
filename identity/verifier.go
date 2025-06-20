package identity

import (
	"context"
	"net/http"

	pe "github.com/marqeta/pr-bot/errors"
)

type Verifier interface {
	Verify(ctx context.Context, r *http.Request) (string, error)
}

type STSVerifier struct {
	HTTPClient *http.Client
	Validator  Validator
	Fetcher    Fetcher
}

func (v *STSVerifier) Verify(ctx context.Context, r *http.Request) (string, error) {
	identity, err := v.Fetcher.FetchCallerIdentity(ctx, r)
	if err != nil {
		return "", err
	}

	if err := v.Validator.ValidateIdentity(identity); err != nil {
		return "", pe.UserError(ctx, "unauthorized identity", err)
	}

	return identity.Arn, nil
}
