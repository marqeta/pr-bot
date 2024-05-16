package rate_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/marqeta/pr-bot/rate"
)

func TestLimiterConfig_Update(t *testing.T) {
	type fields struct {
		Default   rate.Limit
		Overrides map[string]rate.Limit
	}
	tests := []struct {
		name         string
		fields       fields
		wantLimitCfg *rate.LimiterConfig
		wantErr      bool
	}{
		{
			name: "Should parse default fields",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59s",
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Window:    59 * time.Second,
					Value:     10,
					WindowStr: "59s",
				},
			},
			wantErr: false,
		},
		{
			name: "Should parse default fields for window > 1hr",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59h",
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Window:    59 * time.Hour,
					Value:     10,
					WindowStr: "59h",
				},
			},
			wantErr: false,
		},
		{
			name: "Should parse override fields",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o": {
						Value:     110,
						WindowStr: "10s",
					},
					"o1": {
						Value:     120,
						WindowStr: "10m",
					},
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Window:    59 * time.Hour,
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o": {
						Window:    10 * time.Second,
						Value:     110,
						WindowStr: "10s",
					},
					"o1": {
						Window:    10 * time.Minute,
						Value:     120,
						WindowStr: "10m",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should parse override fields for window > 1hr",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o": {
						Value:     110,
						WindowStr: "10h",
					},
					"o1": {
						Value:     120,
						WindowStr: "11h",
					},
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Window:    59 * time.Hour,
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o": {
						Window:    10 * time.Hour,
						Value:     110,
						WindowStr: "10h",
					},
					"o1": {
						Window:    11 * time.Hour,
						Value:     120,
						WindowStr: "11h",
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "Should throw error if default is not present",
			fields:       fields{},
			wantLimitCfg: &rate.LimiterConfig{},
			wantErr:      true,
		},
		{
			name: "Should throw error if default window is not a duration",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59a",
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59a",
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error if override window is not a duration",
			fields: fields{
				Default: rate.Limit{
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o1": {
						Value:     120,
						WindowStr: "10a",
					},
				},
			},
			wantLimitCfg: &rate.LimiterConfig{
				Default: rate.Limit{
					Window:    59 * time.Hour,
					Value:     10,
					WindowStr: "59h",
				},
				Overrides: map[string]rate.Limit{
					"o1": {
						Value:     120,
						WindowStr: "10a",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &rate.LimiterConfig{
				Default:   tt.fields.Default,
				Overrides: tt.fields.Overrides,
			}
			if err := cfg.Update(); (err != nil) != tt.wantErr {
				t.Errorf("LimiterConfig.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantLimitCfg, cfg)
		})
	}
}

func TestLimiterConfig_Get(t *testing.T) {
	type fields struct {
		Default   rate.Limit
		Overrides map[string]rate.Limit
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rate.Limit
	}{
		{
			name: "should return default when no overrides",
			fields: fields{
				Default: limit(1, 1),
			},
			args: args{
				key: "",
			},
			want: limit(1, 1),
		},
		{
			name: "should return default for key without overrides",
			fields: fields{
				Default: limit(1, 1),
				Overrides: map[string]rate.Limit{
					"def": limit(2, 2),
				},
			},
			args: args{
				key: "asd",
			},
			want: limit(1, 1),
		},
		{
			name: "should return override",
			fields: fields{
				Default: limit(1, 1),
				Overrides: map[string]rate.Limit{
					"def": limit(2, 2),
				},
			},
			args: args{
				key: "def",
			},
			want: limit(2, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &rate.LimiterConfig{
				Default:   tt.fields.Default,
				Overrides: tt.fields.Overrides,
			}
			if got := cfg.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LimiterConfig.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func limit(w time.Duration, v int64) rate.Limit {
	return rate.Limit{
		Window:    w,
		Value:     v,
		WindowStr: w.String(),
	}
}
