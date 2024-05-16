package secrets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/marqeta/pr-bot/secrets"
)

var errTimeOut = errors.New("timed out")

func Test_GetSecret(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    *string
		apiErr  error
		wantErr error
	}{
		{
			name:    "should return read secret",
			id:      "key1",
			want:    aws.String("value1"),
			apiErr:  nil,
			wantErr: nil,
		},
		{
			name:    "error when secret does not exists",
			id:      "key1",
			want:    nil,
			apiErr:  nil,
			wantErr: secrets.ErrSecretDoesNotExist,
		},
		{
			name:    "error when secret is empty",
			id:      "key1",
			want:    aws.String(""),
			apiErr:  nil,
			wantErr: secrets.ErrSecretDoesNotExist,
		},
		{
			name:    "error when api call fails",
			id:      "key1",
			want:    nil,
			apiErr:  errTimeOut,
			wantErr: errTimeOut,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := secrets.NewMockAPI(t)
			api.EXPECT().
				GetSecretValue(context.TODO(),
					// #nosec G601
					&secretsmanager.GetSecretValueInput{SecretId: &tt.id}).
				Return(&secretsmanager.GetSecretValueOutput{SecretString: tt.want}, tt.apiErr).
				Once()

			m := secrets.NewManager(api)
			got, err := m.GetSecret(context.TODO(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("manager.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != aws.ToString(tt.want) {
				t.Errorf("manager.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
