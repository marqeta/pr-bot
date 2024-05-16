package evaluation

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/marqeta/pr-bot/metrics"
)

var ErrReportNotFound = errors.New("report not found")

//go:generate mockery --name Manager
type Manager interface {
	NewReportBuilder(ctx context.Context, pr, reqID, deliveryID string) ReportBuilder
	GetReport(ctx context.Context, pr, deliveryID string) (*Report, error)
	StoreReport(ctx context.Context, builder ReportBuilder) error
	ListReports(ctx context.Context, pr string) ([]ReportMetadata, error)
}

type manager struct {
	dao           Dao
	table         string
	ttl           time.Duration
	metrics       metrics.Emitter
	policyVersion string
}

// ListReports implements Manager.
func (m *manager) ListReports(ctx context.Context, pr string) ([]ReportMetadata, error) {
	reports := make([]ReportMetadata, 0)
	keyEx := expression.Key("pr").Equal(expression.Value(pr))
	projEx := expression.NamesList(expression.Name("pr"),
		expression.Name("title"),
		expression.Name("author"),
		expression.Name("request_id"),
		expression.Name("delivery_id"),
		expression.Name("policy_version"),
		expression.Name("outcome"),
		expression.Name("created_at"),
		expression.Name("expire_at"),
		expression.Name("event"),
		expression.Name("action"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).WithProjection(projEx).Build()
	if err != nil {
		return reports, err
	}

	queryPaginator := dynamodb.NewQueryPaginator(m.dao, &dynamodb.QueryInput{
		TableName:                 aws.String(m.table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
	})
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return reports, err
		}
		var reportPage []ReportMetadata
		err = attributevalue.UnmarshalListOfMapsWithOptions(response.Items, &reportPage, useJSONTagForDecoding)
		if err != nil {
			return reports, err
		}
		reports = append(reports, reportPage...)
	}

	if len(reports) >= 0 {
		slices.SortFunc(reports, func(i, j ReportMetadata) int {
			return cmp.Compare(i.CreatedAt, j.CreatedAt)
		})
	}

	return reports, err
}

// GetReport implements Manager.
func (m *manager) GetReport(ctx context.Context, pr, deliveryID string) (*Report, error) {
	key, err := key(pr, deliveryID)
	if err != nil {
		m.emitError(ctx, "GetReport", "MarshalError")
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: &m.table,
	}
	o, err := m.dao.GetItem(context.Background(), input)
	if err != nil {
		m.emitError(ctx, "GetReport", "DDBGetItemError")
		return nil, err
	}

	if o.Item == nil {
		m.emitError(ctx, "GetReport", "ReportNotFound")
		return nil, fmt.Errorf("report for pr: %v delivery_id: %v %w", pr, deliveryID, ErrReportNotFound)
	}

	var report Report
	err = attributevalue.UnmarshalMapWithOptions(o.Item, &report, useJSONTagForDecoding)
	if err != nil {
		m.emitError(ctx, "GetReport", "UnMarshalError")
		return nil, err
	}
	return &report, nil
}

// NewTracker implements Manager.
func (m *manager) NewReportBuilder(_ context.Context, pr, reqID, deliveryID string) ReportBuilder {
	now := time.Now()
	return &reportBuilder{
		report: Report{
			ReportMetadata: ReportMetadata{
				PR:            pr,
				RequestID:     reqID,
				DeliveryID:    deliveryID,
				PolicyVersion: m.policyVersion,
				ExpireAt:      now.Add(m.ttl).Unix(),
				CreatedAt:     now.Unix(),
			},
			Breakdown: make(map[string]Result),
		},
	}
}

// StoreReport implements Manager.
func (m *manager) StoreReport(ctx context.Context, tracker ReportBuilder) error {
	item, err := attributevalue.MarshalMapWithOptions(tracker.GetReport(), useJSONTagEncoding)
	if err != nil {
		m.emitError(ctx, "StoreReport", "MarshalError")
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &m.table,
	}
	_, err = m.dao.PutItem(context.Background(), input)
	if err != nil {
		m.emitError(ctx, "StoreReport", "DDBPutItemError")
		return err
	}
	return nil
}

func NewManager(dao Dao, policyVersion string, ttl time.Duration,
	metrics metrics.Emitter, table string) Manager {
	return &manager{
		dao:           dao,
		ttl:           ttl,
		policyVersion: policyVersion,
		table:         table,
		metrics:       metrics,
	}
}

func key(pr, deliveryID string) (map[string]types.AttributeValue, error) {
	p, err := attributevalue.Marshal(pr)
	if err != nil {
		return nil, err
	}
	id, err := attributevalue.Marshal(deliveryID)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{"pr": p, "delivery_id": id}, nil
}

func useJSONTagEncoding(opts *attributevalue.EncoderOptions) {
	opts.TagKey = "json"
}

func useJSONTagForDecoding(opts *attributevalue.DecoderOptions) {
	opts.TagKey = "json"
}

func (m *manager) emitError(ctx context.Context, name string, errCode string) {
	m.metrics.EmitDist(ctx, "evaluation.manager.error", 1, []string{
		fmt.Sprintf("call:%s", name),
		fmt.Sprintf("code:%s", errCode),
	})
}
