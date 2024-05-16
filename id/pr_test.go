package id_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/marqeta/pr-bot/id"
)

func TestPR_ToTags(t *testing.T) {
	tests := []struct {
		name string
		pr   id.PR
		want []string
	}{
		{
			name: "Should return tags for svc account",
			pr:   toPR("owner", "repo", "owner/repo", 1, "svc-foo"),
			want: toTags("owner", "repo", "owner/repo", 1, "svc-foo", "service-account"),
		},
		{
			name: "Should return tags for user account",
			pr:   toPR("owner", "repo", "owner/repo", 1, "foo"),
			want: toTags("owner", "repo", "owner/repo", 1, "foo", "user"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pr.ToTags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PR.ToTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func toTags(owner, repo, fullname string, number int, author, authorType string) []string {
	return []string{fmt.Sprintf("owner:%s", owner),
		fmt.Sprintf("repo:%s", repo),
		fmt.Sprintf("repoFullName:%s", fullname),
		fmt.Sprintf("pr:%d", number),
		fmt.Sprintf("authorType:%s", authorType),
		fmt.Sprintf("author:%s", author)}
}

func toPR(owner, repo, fullname string, number int, author string) id.PR {
	return id.PR{
		Owner:        owner,
		Repo:         repo,
		Number:       number,
		RepoFullName: fullname,
		Author:       author,
	}
}
