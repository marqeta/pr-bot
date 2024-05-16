package opa_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/marqeta/pr-bot/opa/types"
)

func TestV1_Evaluate(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		module string
		input  *input.Model
	}
	tests := []struct {
		name            string
		args            args
		want            types.Result
		setExpectations func(t *rules.MockRules[bool], r *rules.MockRules[types.Review])
		wantErr         bool
	}{
		{
			name: "return result after evaluating track and review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			want: types.Result{
				Track: true,
				Review: types.Review{
					Type: types.Approve,
					Body: "LGTM",
				},
			},
			setExpectations: func(t *rules.MockRules[bool], r *rules.MockRules[types.Review]) {
				t.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(true, nil)
				r.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(types.Review{
					Type: types.Approve,
					Body: "LGTM",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should not evaluate review when track is false",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			want: types.Result{
				Track: false,
			},
			setExpectations: func(t *rules.MockRules[bool], _ *rules.MockRules[types.Review]) {
				t.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(false, nil)
			},
			wantErr: false,
		},
		{
			name: "should throw error when track evaluation fails",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			want: types.Result{},
			setExpectations: func(t *rules.MockRules[bool], _ *rules.MockRules[types.Review]) {
				t.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(false, randomErr)
			},
			wantErr: true,
		},
		{
			name: "should throw error when review evaluation fails",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			want: types.Result{},
			setExpectations: func(t *rules.MockRules[bool], r *rules.MockRules[types.Review]) {
				t.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(true, nil)
				r.EXPECT().Evaluate(ctx, "ci/module/asd", randomModel()).Return(types.Review{}, randomErr)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := rules.NewMockRules[bool](t)
			review := rules.NewMockRules[types.Review](t)
			v1 := opa.NewV1PolicyFromRules(track, review)
			tt.setExpectations(track, review)
			got, err := v1.Evaluate(ctx, tt.args.module, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("V1.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("V1.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func randomGHE() input.GHE {
	return input.GHE{
		Event:  "random Event",
		Action: "random Action",
		PullRequest: &github.PullRequest{
			Number: aws.Int(259),
			Title:  aws.String("Automatic Dockerfile Image Updater"),
			User: &github.User{
				Login: aws.String("svc-ci-dfiu"),
			},
		},
		Repository: &github.Repository{
			Name:     aws.String("terraform-provider-oci"),
			Owner:    &github.User{Login: aws.String("ci")},
			FullName: aws.String("ci/terraform-provider-oci"),
		},
	}
}

func randomModel() *input.Model {
	ghe := randomGHE()
	return &input.Model{
		Event:       ghe.Event,
		Action:      ghe.Action,
		PullRequest: ghe.PullRequest,
		Repository:  ghe.Repository,
		Plugins: map[string]json.RawMessage{
			"name1": JSONRaw("random message 1"),
		},
	}

}

func JSONRaw(s string) json.RawMessage {
	bytes, _ := json.Marshal(s)
	return json.RawMessage(bytes)
}
