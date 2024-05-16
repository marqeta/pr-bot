package evaluation_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
)

func Test_GetDeliveryID(t *testing.T) {
	ctx := context.TODO()
	id := evaluation.GetDeliveryID(ctx)
	assert.Equal(t, "", id)
	ctx = context.WithValue(ctx,
		evaluation.DeliveryIDKey, "randomID")
	id = evaluation.GetDeliveryID(ctx)
	assert.Equal(t, "randomID", id)
}

func Test_manager_NewReportBuilder(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		ttl           time.Duration
		policyVersion string
	}
	type args struct {
		pr         string
		reqID      string
		deliveryID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   evaluation.Report
	}{
		{
			name: "NewReportBuilder sets fields other than input, outcome and breakdown",
			fields: fields{
				ttl:           1 * time.Hour,
				policyVersion: "version",
			},
			args: args{
				pr:    "pr",
				reqID: "request_id",
			},
			want: evaluation.Report{
				ReportMetadata: evaluation.ReportMetadata{
					PR:            "pr",
					RequestID:     "request_id",
					PolicyVersion: "version",
					ExpireAt:      time.Now().Add(1 * time.Hour).Unix(),
					CreatedAt:     time.Now().Unix(),
				},
				Breakdown: map[string]evaluation.Result{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDao := evaluation.NewMockDao(t)
			m := evaluation.NewManager(mockDao, tt.fields.policyVersion, tt.fields.ttl,
				metrics.NewNoopEmitter(), "table")
			b := m.NewReportBuilder(ctx, tt.args.pr, tt.args.reqID, tt.args.deliveryID)
			if got := b.GetReport(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("manager.NewReportBuilder().GetReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
