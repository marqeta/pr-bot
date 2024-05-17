package rules_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/open-policy-agent/opa/sdk"
)

func TestSchema_Evaluate(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		module string
		input  *input.Model
	}
	tests := []struct {
		name            string
		RuleName        string
		args            args
		setExpectations func(c *client.MockClient)
		want            string
		wantErr         bool
	}{
		{
			name:     "Should return v1 when schema=v1",
			RuleName: "schema",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "schema").Return("ci/module/asd/schema")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/schema",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "v1",
				}, nil)
			},
			want:    "v1",
			wantErr: false,
		},
		{
			name:     "Should return V1 when schema=V1",
			RuleName: "schema",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "schema").Return("ci/module/asd/schema")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/schema",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "V1",
				}, nil)
			},
			want:    "V1",
			wantErr: false,
		},
		{
			name:     "Should return '' when schema is empty",
			RuleName: "schema",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "schema").Return("ci/module/asd/schema")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/schema",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "",
				}, nil)
			},
			want:    "",
			wantErr: false,
		},
		{
			name:     "Should return error when schema returns error",
			RuleName: "schema",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "schema").Return("ci/module/asd/schema")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/schema",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(nil, randomErr)
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewMockClient(t)
			s := rules.NewSchema(tt.RuleName, c)
			tt.setExpectations(c)
			got, err := s.Evaluate(ctx, tt.args.module, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Schema.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
