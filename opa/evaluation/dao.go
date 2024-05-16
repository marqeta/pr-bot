package evaluation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:generate mockery --name Dao  --testonly
type Dao interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}
