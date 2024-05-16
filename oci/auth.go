package oci

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"oras.land/oras-go/v2/registry/remote/auth"
)

var errInvalidToken = fmt.Errorf("invalid authorization token, expected format: <username>:<password>")
var errEmptyAuthorizationData = fmt.Errorf("ecr.GetAuthorizationToken() succeeded but authorization data is empty")

//go:generate mockery --name TokenGetter --testonly
type TokenGetter interface {
	GetAuthorizationToken(ctx context.Context, input *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)
}

type CredentialRetriever interface {
	RetrieveCredential(ctx context.Context) (auth.Credential, error)
}

type ecrCredRetriever struct {
	tokenGetter TokenGetter
}

func NewECRCredRetriever(tokenGeter TokenGetter) CredentialRetriever {
	return &ecrCredRetriever{
		tokenGetter: tokenGeter,
	}
}

func (e *ecrCredRetriever) RetrieveCredential(ctx context.Context) (auth.Credential, error) {
	output, err := e.tokenGetter.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return auth.EmptyCredential, err
	}
	if len(output.AuthorizationData) == 0 {
		return auth.EmptyCredential, errEmptyAuthorizationData
	}
	decoded, err := decode(aws.ToString(output.AuthorizationData[0].AuthorizationToken))
	if err != nil {
		return auth.EmptyCredential, err
	}
	userAndPass := strings.Split(decoded, ":")
	if len(userAndPass) != 2 {
		return auth.EmptyCredential, errInvalidToken
	}
	user := userAndPass[0]
	pass := userAndPass[1]
	return auth.Credential{
		Username: user,
		Password: pass,
	}, nil
}

func decode(token string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
