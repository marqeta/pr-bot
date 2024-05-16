package plugins

import (
	"context"
	"encoding/json"
	"errors"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
)

type branchProtection struct {
	dao gh.API
}

var errNoBaseBranch = errors.New("base branch not found in pull request event payload")

// GetMessage implements input.Plugin.
func (fc *branchProtection) GetInputMsg(ctx context.Context, ghe input.GHE) (json.RawMessage, error) {
	id := ghe.ToID()
	base := ghe.PullRequest.GetBase()
	if base == nil {
		return json.RawMessage{}, errNoBaseBranch
	}
	branch := base.GetRef()
	protection, err := fc.dao.GetBranchProtection(ctx, id, branch)
	if err != nil {
		return json.RawMessage{}, err
	}
	data, err := json.Marshal(protection)
	if err != nil {
		return json.RawMessage{}, err
	}
	return json.RawMessage(data), nil
}

// Name implements input.Plugin.
func (fc *branchProtection) Name() string {
	return "base_branch_protection"
}

func NewBranchProtection(dao gh.API) input.Plugin {
	return &branchProtection{dao: dao}
}
