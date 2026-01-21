package rules_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/open-policy-agent/opa/v1/sdk"
)

func TestTrack_Evaluate(t *testing.T) {
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
		want            bool
		wantErr         bool
	}{
		{
			name:     "Should return true when track=true",
			RuleName: "track",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "track").Return("ci/module/asd/track")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/track",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: true,
				}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "Should return false when track=false",
			RuleName: "track",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "track").Return("ci/module/asd/track")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/track",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(&sdk.DecisionResult{
					ID:     "random decision id",
					Result: false,
				}, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "Should return error when track returns error",
			RuleName: "track",
			args: args{
				module: "ci/module/asd",
				input:  randomModel(),
			},
			setExpectations: func(c *client.MockClient) {
				c.EXPECT().Path("ci/module/asd", "track").Return("ci/module/asd/track")
				c.EXPECT().DecisionID(ctx).Return("random decision id")
				c.EXPECT().Decision(ctx, sdk.DecisionOptions{
					Path:       "ci/module/asd/track",
					Input:      randomModel(),
					DecisionID: "random decision id",
				}).Return(nil, randomErr)
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewMockClient(t)
			tr := rules.NewTrack(tt.RuleName, c)
			tt.setExpectations(c)
			got, err := tr.Evaluate(ctx, tt.args.module, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Track.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Track.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func JSONRaw(s string) json.RawMessage {
	bytes, _ := json.Marshal(s)
	return json.RawMessage(bytes)
}

func randomGHE() input.GHE {
	return input.GHE{
		Event:  "random Event",
		Action: "random Action",
	}
}

func randomModel() *input.Model {
	ghe := randomGHE()
	return &input.Model{
		Event:  ghe.Event,
		Action: ghe.Action,
		Plugins: map[string]json.RawMessage{
			"name1": JSONRaw("random message 1"),
		},
	}

}
