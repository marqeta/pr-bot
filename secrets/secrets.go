package secrets

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var ErrSecretDoesNotExist = errors.New("secret does not exist")

type Manager interface {
	GetSecret(ctx context.Context, id string) (string, error)
}

//go:generate mockery --name API --testonly
type API interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type manager struct {
	api API
}

func NewManager(api API) Manager {
	return manager{
		api: api,
	}
}

func (m manager) GetSecret(ctx context.Context, id string) (string, error) {
	result, err := m.api.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &id})
	if err != nil {
		return "", err
	}

	if result == nil || len(aws.ToString(result.SecretString)) == 0 {
		return "", ErrSecretDoesNotExist
	}

	return aws.ToString(result.SecretString), nil
}
