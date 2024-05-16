package input_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/opa/input"
)

func Test_factory_CreateModel(t *testing.T) {
	ctx := context.TODO()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	tests := []struct {
		name            string
		ghe             input.GHE
		numPlugins      int
		setExpectations func(plugins []*input.MockPlugin)
		want            *input.Model
		wantErr         bool
	}{
		{
			name:       "Should create model with no plugins",
			ghe:        randomGHE(),
			numPlugins: 0,
			setExpectations: func(_ []*input.MockPlugin) {
			},
			want:    toModel(randomGHE(), make(map[string]json.RawMessage)),
			wantErr: false,
		},
		{
			name:       "Should create model with single plugins",
			ghe:        randomGHE(),
			numPlugins: 1,
			setExpectations: func(plugins []*input.MockPlugin) {
				p := plugins[0]
				p.EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message"), nil)
				p.EXPECT().Name().Return("random name")
			},
			want: toModel(randomGHE(),
				map[string]json.RawMessage{"random name": JSONRaw("random message")}),
			wantErr: false,
		},
		{
			name:       "Should create model with multiple plugins",
			ghe:        randomGHE(),
			numPlugins: 3,
			setExpectations: func(plugins []*input.MockPlugin) {
				plugins[0].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message 0"), nil)
				plugins[0].EXPECT().Name().Return("name0")
				plugins[1].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message 1"), nil)
				plugins[1].EXPECT().Name().Return("name1")
				plugins[2].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message 2"), nil)
				plugins[2].EXPECT().Name().Return("name2")
			},
			want: toModel(randomGHE(),
				map[string]json.RawMessage{
					"name0": JSONRaw("random message 0"),
					"name1": JSONRaw("random message 1"),
					"name2": JSONRaw("random message 2"),
				}),
			wantErr: false,
		},
		{
			name:       "Should skip msg from plugins with error",
			ghe:        randomGHE(),
			numPlugins: 3,
			setExpectations: func(plugins []*input.MockPlugin) {
				plugins[0].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message 0"), nil)
				plugins[0].EXPECT().Name().Return("name0")
				plugins[1].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(nil, randomErr)
				plugins[1].EXPECT().Name().Return("name1")
				plugins[2].EXPECT().GetInputMsg(ctx, randomGHE()).
					Return(JSONRaw("random message 2"), nil)
				plugins[2].EXPECT().Name().Return("name2")
			},
			want: toModel(randomGHE(),
				map[string]json.RawMessage{
					"name0": JSONRaw("random message 0"),
					"name2": JSONRaw("random message 2")},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlugins := make([]*input.MockPlugin, tt.numPlugins)
			plugins := make([]input.Plugin, tt.numPlugins)
			for i := 0; i < tt.numPlugins; i++ {
				p := input.NewMockPlugin(t)
				plugins[i] = p
				mockPlugins[i] = p
			}
			tt.setExpectations(mockPlugins)
			f := input.NewFactory(plugins...)
			got, err := f.CreateModel(ctx, tt.ghe)
			if (err != nil) != tt.wantErr {
				t.Errorf("factory.CreateModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("factory.CreateModel() = %+v, want %+v", got, tt.want)
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
		PullRequest: &github.PullRequest{
			Number: aws.Int(1),
			NodeID: aws.String("random NodeID"),
			User: &github.User{
				Login: aws.String("random Login"),
			},
		},
		Repository: &github.Repository{
			Name:     aws.String("random Name"),
			Owner:    &github.User{Login: aws.String("random Login")},
			FullName: aws.String("random FullName"),
		},
	}
}

func toModel(ghe input.GHE, plugins map[string]json.RawMessage) *input.Model {
	return &input.Model{
		Event:        ghe.Event,
		Action:       ghe.Action,
		PullRequest:  ghe.PullRequest,
		Repository:   ghe.Repository,
		Organization: ghe.Organization,
		Plugins:      plugins,
	}
}
