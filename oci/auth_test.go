package oci_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/marqeta/pr-bot/oci"
	"oras.land/oras-go/v2/registry/remote/auth"
)

func Test_ecrCredRetriever_RetrieveCredential(t *testing.T) {
	ctx := context.TODO()

	//nolint:goerr113
	err := fmt.Errorf("random")
	tests := []struct {
		name            string
		setExpectations func(g *oci.MockTokenGetter)
		want            auth.Credential
		wantErr         bool
	}{
		{
			name: "Should return credential",
			setExpectations: func(g *oci.MockTokenGetter) {
				g.EXPECT().
					GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{}).
					Return(mockOutput("user", "pass"), nil)
			},
			want: auth.Credential{
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "Should return error when token is invalid",
			setExpectations: func(g *oci.MockTokenGetter) {
				g.EXPECT().
					GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{}).
					Return(InvalidTokenOutput("user", "pass"), nil)
			},
			want:    auth.EmptyCredential,
			wantErr: true,
		},
		{
			name: "Should return error when Authorization Data is empty is empty",
			setExpectations: func(g *oci.MockTokenGetter) {
				g.EXPECT().
					GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{}).
					Return(&ecr.GetAuthorizationTokenOutput{}, nil)
			},
			want:    auth.EmptyCredential,
			wantErr: true,
		},
		{
			name: "Should return error when GetAuthorizationToken fails",
			setExpectations: func(g *oci.MockTokenGetter) {
				g.EXPECT().
					GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{}).
					Return(&ecr.GetAuthorizationTokenOutput{}, err)
			},
			want:    auth.EmptyCredential,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := oci.NewMockTokenGetter(t)
			e := oci.NewECRCredRetriever(g)
			tt.setExpectations(g)
			got, err := e.RetrieveCredential(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("ecrCredRetriever.RetrieveCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ecrCredRetriever.RetrieveCredential() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockOutput(user, pass string) *ecr.GetAuthorizationTokenOutput {
	token := fmt.Sprintf("%s:%s", user, pass)
	encoded := base64.StdEncoding.EncodeToString([]byte(token))
	return &ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []types.AuthorizationData{{
			AuthorizationToken: &encoded,
		}},
	}
}

func InvalidTokenOutput(user, pass string) *ecr.GetAuthorizationTokenOutput {
	token := fmt.Sprintf("%s_%s", user, pass)
	encoded := base64.StdEncoding.EncodeToString([]byte(token))
	return &ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []types.AuthorizationData{{
			AuthorizationToken: &encoded,
		}},
	}
}
