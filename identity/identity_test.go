package identity_test

import (
	"context"
	"errors"
	"testing"

	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/identity"
	"github.com/stretchr/testify/assert"
)

// Wrap the validator to take a config
type testValidator struct {
	Config *prbot.Config
}

func (v *testValidator) ValidateIdentity(ctx context.Context, identity *identity.CallerIdentity) error {
	if identity == nil ||
		!stringInSlice(identity.Arn, v.Config.Identity.AllowedCallerArns) ||
		!stringInSlice(identity.Account, v.Config.Identity.AllowedCallerAccounts) {
		return errors.New("identity missing or invalid")
	}
	return nil
}

func stringInSlice(target string, list []string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func TestAllowAllValidator_ValidateIdentity(t *testing.T) {
	ctx := context.Background()

	testConfig := &prbot.Config{
		Identity: struct {
			AllowedCallerArns     []string `yaml:"AllowedCallerArns" env:"ALLOWED_CALLER_ARNS"`
			AllowedCallerAccounts []string `yaml:"AllowedCallerAccounts" env:"ALLOWED_CALLER_ACCOUNTS"`
		}{
			AllowedCallerArns:     []string{"arn:aws:iam::123456789012:role/MyRole"},
			AllowedCallerAccounts: []string{"123456789012"},
		},
	}

	validator := &testValidator{Config: testConfig}

	tests := []struct {
		name    string
		input   *identity.CallerIdentity
		wantErr error
	}{
		{
			name:    "valid identity",
			input:   &identity.CallerIdentity{Arn: "arn:aws:iam::123456789012:role/MyRole", Account: "123456789012"},
			wantErr: nil,
		},
		{
			name:    "nil identity",
			input:   nil,
			wantErr: errors.New("identity missing or invalid"),
		},
		{
			name:    "invalid arn",
			input:   &identity.CallerIdentity{Arn: "arn:aws:iam::000000000000:role/InvalidRole", Account: "123456789012"},
			wantErr: errors.New("identity missing or invalid"),
		},
		{
			name:    "invalid account",
			input:   &identity.CallerIdentity{Arn: "arn:aws:iam::123456789012:role/MyRole", Account: "000000000000"},
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
