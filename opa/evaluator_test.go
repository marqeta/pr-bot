package opa_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
)

func Test_evaluator_Evaluate(t *testing.T) {
	ctx := context.WithValue(context.TODO(), middleware.RequestIDKey, "request_id")
	ctx = context.WithValue(ctx, evaluation.DeliveryIDKey, "delivery_id")
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		ghe input.GHE
	}
	tests := []struct {
		name            string
		modules         []string
		args            args
		setExpectations func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder)
		want            types.Result
		wantErr         bool
	}{
		{
			name:    "should return result with one module",
			modules: []string{"m1"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(approve(), nil))
				b.EXPECT().SetOutcome(evalResult(approve(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    approve(),
			wantErr: false,
		},
		{
			name:    "should return req_changes when modules returns req_changes as highest priority",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(comment(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(comment(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(reqChanges(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(reqChanges(), nil))
				b.EXPECT().SetOutcome(evalResult(reqChanges(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    reqChanges(),
			wantErr: false,
		},
		{
			name:    "should return comment when modules returns comment as highest priority",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(comment(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(comment(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(notTrack(), nil))
				b.EXPECT().SetOutcome(evalResult(comment(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    comment(),
			wantErr: false,
		},
		{
			name:    "should return approve when modules returns approve as highest priority",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(notTrack(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(skip(), nil))
				b.EXPECT().SetOutcome(evalResult(approve(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    approve(),
			wantErr: false,
		},
		{
			name:    "should return skip when modules returns skip as highest priority",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(notTrack(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(skip(), nil))
				b.EXPECT().SetOutcome(evalResult(skip(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    skip(),
			wantErr: false,
		},
		{
			name:    "should return not track when all modules returns not track",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(notTrack(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(notTrack(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(notTrack(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(notTrack(), nil))
				b.EXPECT().SetOutcome(evalResult(notTrack(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    notTrack(),
			wantErr: false,
		},
		{
			name:    "should evaluate next modules when req_changes is returned",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(reqChanges(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(reqChanges(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(skip(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(skip(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(notTrack(), nil)
				b.EXPECT().AddModuleResult("m4", evalResult(notTrack(), nil))
				b.EXPECT().SetOutcome(evalResult(reqChanges(), nil))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    reqChanges(),
			wantErr: false,
		},
		{
			name:    "should throw error when create model fails",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(_ *opa.MockPolicy, f *input.MockFactory, _ *evaluation.MockManager, _ *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(&input.Model{}, randomErr)
			},
			want:    notTrack(),
			wantErr: true,
		},
		{
			name:    "should return error when single module eval fails",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m1", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m2", randomModel()).Return(approve(), nil)
				b.EXPECT().AddModuleResult("m2", evalResult(approve(), nil))
				p.EXPECT().Evaluate(ctx, "m3", randomModel()).Return(reqChanges(), nil)
				b.EXPECT().AddModuleResult("m3", evalResult(reqChanges(), nil))
				p.EXPECT().Evaluate(ctx, "m4", randomModel()).Return(notTrack(), randomErr)
				b.EXPECT().AddModuleResult("m4", evalResult(notTrack(), randomErr))
				b.EXPECT().SetOutcome(evalResult(notTrack(), randomErr))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    notTrack(),
			wantErr: true,
		},
		{
			name:    "can skip other modules when first module eval fails",
			modules: []string{"m1", "m2", "m3", "m4"},
			args: args{
				ghe: randomGHE(),
			},
			setExpectations: func(p *opa.MockPolicy, f *input.MockFactory, m *evaluation.MockManager, b *evaluation.MockReportBuilder) {
				f.EXPECT().CreateModel(ctx, randomGHE()).Return(randomModel(), nil)
				m.EXPECT().NewReportBuilder(ctx, "ci/terraform-provider-oci/259", "request_id", "delivery_id").Return(b)
				b.EXPECT().SetInput(randomModel())
				p.EXPECT().Evaluate(ctx, "m1", randomModel()).Return(notTrack(), randomErr)
				b.EXPECT().AddModuleResult("m1", evalResult(notTrack(), randomErr))
				b.EXPECT().SetOutcome(evalResult(notTrack(), randomErr))
				m.EXPECT().StoreReport(ctx, b).Return(nil)
			},
			want:    notTrack(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := input.NewMockFactory(t)
			p := opa.NewMockPolicy(t)
			m := evaluation.NewMockManager(t)
			b := evaluation.NewMockReportBuilder(t)

			e := opa.NewEvaluator(tt.modules, p, f, m)
			tt.setExpectations(p, f, m, b)
			got, err := e.Evaluate(ctx, tt.args.ghe)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evaluator.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func evalResult(r types.Result, e error) evaluation.Result {
	return evaluation.Result{
		Result: r,
		Err:    e,
	}
}

func approve() types.Result {
	return types.Result{
		Track: true,
		Review: types.Review{
			Type: types.Approve,
			Body: "LGTM",
		},
	}
}

func reqChanges() types.Result {
	return types.Result{
		Track: true,
		Review: types.Review{
			Type: types.RequestChanges,
			Body: "need unit tests",
		},
	}
}

func comment() types.Result {
	return types.Result{
		Track: true,
		Review: types.Review{
			Type: types.Comment,
			Body: "nit spelling",
		},
	}
}

func skip() types.Result {
	return types.Result{
		Track: true,
		Review: types.Review{
			Type: types.Skip,
		},
	}
}

func notTrack() types.Result {
	return types.Result{}
}
