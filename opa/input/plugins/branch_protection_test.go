package plugins_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/input/plugins"
)

func TestBranchProtection_GetInputMsg(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	type args struct {
		ghe            input.GHE
		setExpectaions func(d *gh.MockAPI)
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Should return branch protection rule for base branch",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().GetBranchProtection(ctx, randomGHE().ToID(), "random Ref").
						Return(randomBranchProtection(), nil)
				},
			},
			want:    toJSON(t, randomBranchProtection()),
			wantErr: false,
		},
		{
			name: "Should return error when dao returns error",
			args: args{
				ghe: randomGHE(),
				setExpectaions: func(d *gh.MockAPI) {
					d.EXPECT().GetBranchProtection(ctx, randomGHE().ToID(), "random Ref").
						Return(nil, randomErr)
				},
			},
			want:    json.RawMessage([]byte{}),
			wantErr: true,
		},
		{
			name: "Should return error when base branch is empty",
			args: args{
				ghe:            emptyBaseBranch(),
				setExpectaions: func(_ *gh.MockAPI) {},
			},
			want:    json.RawMessage([]byte{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := gh.NewMockAPI(t)
			bbp := plugins.NewBranchProtection(dao)
			tt.args.setExpectaions(dao)
			got, err := bbp.GetInputMsg(ctx, tt.args.ghe)
			if (err != nil) != tt.wantErr {
				t.Errorf("BranchProtection.GetInputMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BranchProtection.GetInputMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func randomBranchProtection() *github.Protection {
	return &github.Protection{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   false,
			Contexts: []string{"random context"},
			Checks: []*github.RequiredStatusCheck{{
				Context: "random context",
				AppID:   aws.Int64(123),
			}},
		},
	}
}

func emptyBaseBranch() input.GHE {
	ghe := randomGHE()
	ghe.PullRequest.Base = nil
	return ghe
}
