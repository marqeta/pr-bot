package input

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/id"
)

type GHE struct {
	Event        string               `json:"event"`
	Action       string               `json:"action"`
	PullRequest  *github.PullRequest  `json:"pull_request,omitempty"`
	Repository   *github.Repository   `json:"repository,omitempty"`
	Organization *github.Organization `json:"organization,omitempty"`
}

type Model struct {
	Event        string                     `json:"event"`
	Action       string                     `json:"action"`
	PullRequest  *github.PullRequest        `json:"pull_request,omitempty"`
	Repository   *github.Repository         `json:"repository,omitempty"`
	Organization *github.Organization       `json:"organization,omitempty"`
	Plugins      map[string]json.RawMessage `json:"plugins"`
}

func (ghe GHE) ToID() id.PR {
	return id.PR{
		Owner:        aws.ToString(ghe.Repository.Owner.Login),
		Repo:         aws.ToString(ghe.Repository.Name),
		Number:       aws.ToInt(ghe.PullRequest.Number),
		NodeID:       aws.ToString(ghe.PullRequest.NodeID),
		RepoFullName: aws.ToString(ghe.Repository.FullName),
		Author:       aws.ToString(ghe.PullRequest.User.Login),
	}
}
