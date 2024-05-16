package input

import (
	"context"
	"encoding/json"

	"github.com/go-chi/httplog"
)

//go:generate mockery --name Factory
type Factory interface {
	CreateModel(ctx context.Context, ghe GHE) (*Model, error)
}

type factory struct {
	plugins []Plugin
}

// CreateModel implements Factory.
func (f *factory) CreateModel(ctx context.Context, ghe GHE) (*Model, error) {

	oplog := httplog.LogEntry(ctx)
	model := Model{
		Event:        ghe.Event,
		Action:       ghe.Action,
		PullRequest:  ghe.PullRequest,
		Repository:   ghe.Repository,
		Organization: ghe.Organization,
		Plugins:      make(map[string]json.RawMessage),
	}

	for _, plugin := range f.plugins {
		name := plugin.Name()
		msg, err := plugin.GetInputMsg(ctx, ghe)
		if err != nil {
			// TODO add metric
			oplog.Err(err).Msgf("failed to get message for plugin %s", name)
			continue
		}
		model.Plugins[name] = msg
	}

	return &model, nil
}

func NewFactory(plugins ...Plugin) Factory {
	return &factory{plugins: plugins}
}
