package rate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/rate"
)

func Test_facade_ShouldThrottle(t *testing.T) {

	//nolint:goerr113
	randErr := errors.New("random error")
	ctx := context.Background()
	type fields struct {
		throttlers []rate.Throttler
	}
	type args struct {
		ID id.PR
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "facade => error when one throttler => error",
			fields: fields{
				throttlers: []rate.Throttler{
					rate.NewMockThrottler(nil), rate.NewMockThrottler(nil), rate.NewMockThrottler(randErr),
				},
			},
			args: args{
				ID: id.PR{},
			},
			wantErr: true,
		},
		{
			name: "facade => error when all throttler => error",
			fields: fields{
				throttlers: []rate.Throttler{
					rate.NewMockThrottler(randErr), rate.NewMockThrottler(randErr), rate.NewMockThrottler(randErr),
				},
			},
			args: args{
				ID: id.PR{},
			},
			wantErr: true,
		},
		{
			name: "facade => nil when all throttler => nil",
			fields: fields{
				throttlers: []rate.Throttler{
					rate.NewMockThrottler(nil), rate.NewMockThrottler(nil), rate.NewMockThrottler(nil),
				},
			},
			args: args{
				ID: id.PR{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := rate.NewFacade(metrics.NewNoopEmitter(), tt.fields.throttlers...)
			if err := f.ShouldThrottle(ctx, tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("facade.ShouldThrottle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
