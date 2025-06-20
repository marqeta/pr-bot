package identity_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/marqeta/pr-bot/identity"
	"github.com/stretchr/testify/assert"
)

type mockFetcher struct {
	identity *identity.CallerIdentity
	err      error
}

func (m *mockFetcher) FetchCallerIdentity(ctx context.Context, r *http.Request) (*identity.CallerIdentity, error) {
	return m.identity, m.err
}

type mockValidator struct {
	err error
}

func (m *mockValidator) ValidateIdentity(ctx context.Context, identity *identity.CallerIdentity) error {
	return m.err
}

func TestSTSVerifier_Verify(t *testing.T) {
	ctx := context.TODO()
	req, _ := http.NewRequest("GET", "http://localhost", nil)

	t.Run("valid identity", func(t *testing.T) {
		v := identity.STSVerifier{
			Fetcher:   &mockFetcher{identity: &identity.CallerIdentity{Arn: "arn:aws:iam::abc"}, err: nil},
			Validator: &mockValidator{err: nil},
		}
		arn, err := v.Verify(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "arn:aws:iam::abc", arn)
	})

	t.Run("fetch error", func(t *testing.T) {
		v := identity.STSVerifier{
			Fetcher:   &mockFetcher{identity: nil, err: errors.New("fetch failed")},
			Validator: &mockValidator{err: nil},
		}
		arn, err := v.Verify(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, arn)
	})

	t.Run("validation error", func(t *testing.T) {
		v := identity.STSVerifier{
			Fetcher:   &mockFetcher{identity: &identity.CallerIdentity{Arn: "bad"}, err: nil},
			Validator: &mockValidator{err: errors.New("unauthorized")},
		}
		arn, err := v.Verify(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized identity")
		assert.Empty(t, arn)
	})
}
