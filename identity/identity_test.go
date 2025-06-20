package identity_test

import (
	"context"
	"errors"
	"testing"

	"github.com/marqeta/pr-bot/identity"
	"github.com/stretchr/testify/assert"
)

func TestAllowAllValidator_ValidateIdentity(t *testing.T) {
	validator := &identity.AllowAllValidator{}
	ctx := context.Background()

	tests := []struct {
		name    string
		input   *identity.CallerIdentity
		wantErr error
	}{
		{
			name:    "valid identity",
			input:   &identity.CallerIdentity{Arn: "arn:aws:iam::123456789012:role/MyRole"},
			wantErr: nil,
		},
		{
			name:  "nil identity",
			input: nil,
			//nolint:goerr113
			wantErr: errors.New("identity missing or invalid"),
		},
		{
			name:  "empty ARN",
			input: &identity.CallerIdentity{Arn: ""},
			//nolint:goerr113
			wantErr: errors.New("identity missing or invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateIdentity(ctx, tt.input)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
