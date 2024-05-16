package opa_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/marqeta/pr-bot/opa/types"
)

func Test_versionedPolicy_Evaluate(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		module string
		ip     *input.Model
	}
	tests := []struct {
		name            string
		args            args
		setupVersions   func(t *testing.T) map[string]*opa.MockPolicy
		setExpectations func(s *rules.MockRules[string], versions map[string]*opa.MockPolicy)
		want            types.Result
		wantErr         bool
	}{
		{
			name: "Should call policy with version v1",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], versions map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("v1", nil)
				p := versions["v1"]
				p.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(approve(), nil)
			},
			want:    approve(),
			wantErr: false,
		},
		{
			name: "Should call policy with version v2",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
					"v2": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], versions map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("v2", nil)
				p := versions["v2"]
				p.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(approve(), nil)
			},
			want:    approve(),
			wantErr: false,
		},
		{
			name: "Should call policy with version V1 (uppercase)",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], versions map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("V1", nil)
				p := versions["v1"]
				p.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(approve(), nil)
			},
			want:    approve(),
			wantErr: false,
		},
		{
			name: "Should throw error when version not found",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], _ map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("V10", nil)
			},
			want:    types.Result{},
			wantErr: true,
		},
		{
			name: "Should throw error when schema evaluation fails",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], _ map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("", randomErr)
			},
			want:    types.Result{},
			wantErr: true,
		},
		{
			name: "Should throw error when policy evaluation fails",
			args: args{
				module: "ci/module/asd",
				ip:     randomModel(),
			},
			setupVersions: func(t *testing.T) map[string]*opa.MockPolicy {
				return map[string]*opa.MockPolicy{
					"v1": opa.NewMockPolicy(t),
				}
			},
			setExpectations: func(schema *rules.MockRules[string], versions map[string]*opa.MockPolicy) {
				schema.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return("v1", nil)
				p := versions["v1"]
				p.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(types.Result{}, randomErr)
			},
			want:    types.Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := rules.NewMockRules[string](t)
			mockVersions := tt.setupVersions(t)
			versions := make(map[string]opa.Policy)
			for k, v := range mockVersions {
				versions[k] = v
			}
			tt.setExpectations(schema, mockVersions)
			p := opa.NewVersionedPolicyFromRules(versions, schema)
			got, err := p.Evaluate(ctx, tt.args.module, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("versionedPolicy.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("versionedPolicy.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
