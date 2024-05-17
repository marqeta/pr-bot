package rules_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/marqeta/pr-bot/opa/types"
	"github.com/open-policy-agent/opa/sdk"
)

func TestReview_Evaluate(t *testing.T) {
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
		want            types.Review
		wantErr         bool
	}{
		{
			name:     "Should return APPROVE",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.Approve.String(),
				}, nil)
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/body",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "LGTM",
				}, nil)
			},
			want: types.Review{
				Type: types.Approve,
				Body: "LGTM",
			},
			wantErr: false,
		},
		{
			name:     "Should return COMMENT",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.Comment.String(),
				}, nil)
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/body",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "random comment",
				}, nil)
			},
			want: types.Review{
				Type: types.Comment,
				Body: "random comment",
			},
			wantErr: false,
		},
		{
			name:     "Should return REQUEST_CHANGES",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.RequestChanges.String(),
				}, nil)
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/body",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "random comment to request changes",
				}, nil)
			},
			want: types.Review{
				Type: types.RequestChanges,
				Body: "random comment to request changes",
			},
			wantErr: false,
		},
		{
			name:     "Should return SKIP",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.Skip.String(),
				}, nil)
			},
			want: types.Review{
				Type: types.Skip,
			},
			wantErr: false,
		},
		{
			name:     "Should return SKIP",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.Skip.String(),
				}, nil)
			},
			want: types.Review{
				Type: types.Skip,
			},
			wantErr: false,
		},
		{
			name:     "Should return error when review.type returns error",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "",
				}, randomErr)
			},
			want:    types.Review{},
			wantErr: true,
		},
		{
			name:     "Should return error when review.body returns error",
			RuleName: "review",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "review").Return("ci/module/asd/review")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/type",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: types.RequestChanges.String(),
				}, nil)
				c.EXPECT().Decision(ctx,
					sdk.DecisionOptions{
						Path:       "ci/module/asd/review/body",
						Input:      randomModel(),
						DecisionID: "random decision id",
					}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: "",
				}, randomErr)
			},
			want:    types.Review{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewMockClient(t)
			r := rules.NewReview(tt.RuleName, c)
			tt.setExpectations(c)
			got, err := r.Evaluate(ctx, tt.args.module, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Review.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Review.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
