package datastore_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/marqeta/pr-bot/datastore"
	"github.com/marqeta/pr-bot/metrics"
)

func Test_dynamo_GetPayload(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		m       *datastore.Metadata
		want    json.RawMessage
		wantErr bool
	}{
		{
			name:    "should get payload",
			m:       randomMetadata(),
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := datastore.NewDynamoDao(nil, metrics.NewNoopEmitter())
			got, err := d.GetPayload(ctx, tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("dynamo.GetPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dynamo.GetPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dynamo_StorePayload(t *testing.T) {
	ctx := context.Background()
	type args struct {
		m       *datastore.Metadata
		payload json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should store payload",
			args: args{
				m:       randomMetadata(),
				payload: json.RawMessage{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := datastore.NewDynamoDao(nil, metrics.NewNoopEmitter())
			if err := d.StorePayload(ctx, tt.args.m, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("dynamo.StorePayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDynamoDao(t *testing.T) {
	if got := datastore.NewDynamoDao(nil, metrics.NewNoopEmitter()); got == nil {
		t.Errorf("NewDynamoDao() = %v, want Dao", got)
	}
}
