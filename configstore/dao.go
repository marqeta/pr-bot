package configstore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/marqeta/pr-bot/metrics"
)

//go:generate mockery --name Dao  --testonly
type Dao[T DynamicConfig] interface {
	GetItem(key string, table string) (T, error)
}

type dynamo[T DynamicConfig] struct {
	client  *dynamodb.Client
	metrics metrics.Emitter
	zero    T
}

// GetItem implements Dao
func (d *dynamo[T]) GetItem(name string, table string) (T, error) {
	key, err := getKey(name)
	if err != nil {
		emitError(d.metrics, name, "MarshalError")
		return d.zero, err
	}
	input := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      &table,
		ConsistentRead: aws.Bool(true),
	}
	o, err := d.client.GetItem(context.Background(), input)
	if err != nil {
		emitError(d.metrics, name, "DDBGetItemError")
		return d.zero, err
	}

	if o.Item == nil {
		emitError(d.metrics, name, "ConfigNotFound")
		//nolint:goerr113
		return d.zero, fmt.Errorf("%v config not found in ddb", name)
	}

	var cfg T
	err = attributevalue.UnmarshalMap(o.Item, &cfg)
	if err != nil {
		emitError(d.metrics, name, "UnMarshalError")
		return d.zero, err
	}
	return cfg, nil
}

func NewDynamoDao[T DynamicConfig](client *dynamodb.Client, m metrics.Emitter) Dao[T] {
	return &dynamo[T]{
		client:  client,
		metrics: m,
	}
}

func getKey(name string) (map[string]types.AttributeValue, error) {
	n, err := attributevalue.Marshal(name)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{"name": n}, nil
}
