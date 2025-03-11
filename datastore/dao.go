package datastore

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/httplog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marqeta/pr-bot/metrics"
)

//go:generate mockery --name Dao
type Dao interface {
	GetPayload(ctx context.Context, m *Metadata) (json.RawMessage, error)
	StorePayload(ctx context.Context, m *Metadata, payload json.RawMessage) error
	ToMetadata(ctx context.Context, r *http.Request) (*Metadata, error)
}

type dynamo struct {
	client  *dynamodb.Client
	metrics metrics.Emitter
}

// GetPayload implements Dao
func (d *dynamo) GetPayload(ctx context.Context, m *Metadata) (json.RawMessage, error) {
	oplog := httplog.LogEntry(ctx)
	oplog.Info().Interface("metadata", m).Msg("getting payload")
	return nil, nil
}

// StorePayload implements Dao
func (d *dynamo) StorePayload(ctx context.Context, m *Metadata, payload json.RawMessage) error {
	oplog := httplog.LogEntry(ctx)
	oplog.Info().Interface("metadata", m).Msg("storing payload")
	return nil
}

// ToMetadata implements Dao
func (d *dynamo) ToMetadata(ctx context.Context, r *http.Request) (*Metadata, error) {
	return ToMetadata(r)
}

func NewDynamoDao(client *dynamodb.Client, m metrics.Emitter) Dao {
	return &dynamo{
		client:  client,
		metrics: m,
	}
}
