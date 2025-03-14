package id

import (
	"fmt"
	"strings"
)

type PR struct {
	Owner        string `json:"owner"`
	Repo         string `json:"repo"`
	Number       int    `json:"number"`
	NodeID       string `json:"node_id,omitempty"`
	RepoFullName string `json:"repo_full_name,omitempty"`
	Author       string `json:"author,omitempty"`
	URL          string `json:"url,omitempty"`
}

func (pr PR) ToTags() []string {
	acctType := "user"
	if strings.HasPrefix(pr.Author, "svc-") {
		acctType = "service-account"
	}
	return []string{fmt.Sprintf("owner:%s", pr.Owner),
		fmt.Sprintf("repo:%s", pr.Repo),
		fmt.Sprintf("repoFullName:%s", pr.RepoFullName),
		fmt.Sprintf("pr:%d", pr.Number),
		fmt.Sprintf("authorType:%s", acctType),
		fmt.Sprintf("author:%s", pr.Author)}
}
