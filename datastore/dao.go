package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/httplog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/marqeta/pr-bot/metrics"
)

type record struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
	Metadata
	Payload   string `json:"payload"`
	ExpireAt  int64  `json:"expire_at"`
	CreatedAt int64  `json:"created_at"`
}

//go:generate mockery --name Dao
type Dao interface {
	GetPayload(ctx context.Context, m *Metadata) (json.RawMessage, error)
	StorePayload(ctx context.Context, m *Metadata, payload json.RawMessage) error
	ToMetadata(ctx context.Context, r *http.Request) (*Metadata, error)
}

type dynamo struct {
	client  *dynamodb.Client
	metrics metrics.Emitter
	ttl     time.Duration
	table   string
}

// GetPayload implements Dao
func (d *dynamo) GetPayload(ctx context.Context, m *Metadata) (json.RawMessage, error) {
	oplog := httplog.LogEntry(ctx)
	oplog.Info().Interface("metadata", m).Msg("getting payload")
	key, err := keyAttributeValue(m)
	if err != nil {
		d.emitError(ctx, "GetPayload", "MarshalError")
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      &d.table,
		ConsistentRead: aws.Bool(true),
	}
	resp, err := d.client.GetItem(ctx, input)
	if err != nil {
		d.emitError(ctx, "GetPayload", "DDBGetItemError")
		return nil, err
	}

	if resp.Item == nil {
		oplog.Info().Interface("metadata", m).Msg("item not found")
		return nil, nil
	}
	var record record
	err = attributevalue.UnmarshalMapWithOptions(resp.Item, &record, useJSONTagForDecoding)
	if err != nil {
		d.emitError(ctx, "GetPayload", "UnMarshalError")
		return nil, err
	}
	return []byte(record.Payload), nil
}

// StorePayload implements Dao
func (d *dynamo) StorePayload(ctx context.Context, m *Metadata, payload json.RawMessage) error {
	oplog := httplog.LogEntry(ctx)
	oplog.Info().Interface("metadata", m).Msg("storing payload")
	record := d.toRecord(m, payload)

	item, err := attributevalue.MarshalMapWithOptions(record, useJSONTagEncoding)
	if err != nil {
		d.emitError(ctx, "StorePayload", "MarshalError")
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &d.table,
	}
	_, err = d.client.PutItem(ctx, input)
	if err != nil {
		d.emitError(ctx, "StorePayload", "DDBPutItemError")
		return err
	}
	return nil
}

// ToMetadata implements Dao
func (d *dynamo) ToMetadata(_ context.Context, r *http.Request) (*Metadata, error) {
	return ToMetadata(r)
}

func NewDynamoDao(client *dynamodb.Client, m metrics.Emitter, ttl time.Duration, table string) Dao {
	return &dynamo{
		client:  client,
		metrics: m,
		ttl:     ttl,
		table:   table,
	}
}

func useJSONTagEncoding(opts *attributevalue.EncoderOptions) {
	opts.TagKey = "json"
}

func useJSONTagForDecoding(opts *attributevalue.DecoderOptions) {
	opts.TagKey = "json"
}

func keyAttributeValue(m *Metadata) (map[string]types.AttributeValue, error) {
	pk, sk := keys(m)
	pkm, err := attributevalue.Marshal(pk)
	if err != nil {
		return nil, err
	}
	skm, err := attributevalue.Marshal(sk)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{"pk": pkm, "sk": skm}, nil
}

func keys(m *Metadata) (string, string) {
	pk := fmt.Sprintf("%s/%s/%d", m.Owner, m.Repo, m.Number)
	sk := fmt.Sprintf("%s/%s/%s/%s", m.Service, m.Job, m.Head, m.Base)
	return pk, sk
}

func (d dynamo) toRecord(m *Metadata, payload json.RawMessage) *record {
	now := time.Now()
	pk, sk := keys(m)
	r := &record{
		PK:        pk,
		SK:        sk,
		Metadata:  *m,
		Payload:   string(payload),
		ExpireAt:  now.Add(d.ttl).Unix(),
		CreatedAt: now.Unix(),
	}
	return r
}

func (d *dynamo) emitError(ctx context.Context, name string, errCode string) {
	d.metrics.EmitDist(ctx, "datastore.dao.error", 1, []string{
		fmt.Sprintf("call:%s", name),
		fmt.Sprintf("code:%s", errCode),
	})
}
