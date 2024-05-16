package id

import (
	"fmt"
	"strings"
)

type PR struct {
	Owner        string
	Repo         string
	Number       int
	NodeID       string
	RepoFullName string
	Author       string
	URL          string
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
