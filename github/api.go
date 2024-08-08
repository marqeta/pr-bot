package github

import (
	"context"

	"github.com/google/go-github/v50/github"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/id"
	"github.com/shurcooL/githubv4"
)

const (
	Approve          = "APPROVE"
	RequestChanges   = "REQUEST_CHANGES"
	Comment          = "COMMENT"
	LGTM             = "LGTM! :rocket: :tada: :100:"
	ApprovalTemplate = `
<details>
<summary>
%v
<br/>
More details on PR bot policy evaluation: <a href="%v">link</a>
</summary>

~~~json
%v
~~~

</details>
`

	ErrorTemplate = `
<details>
	<summary>:warning: %v :warning: </summary>

~~~json
%v
~~~

</details>
`
)

//go:generate mockery --name API
type API interface {
	ListReviews(ctx context.Context, id id.PR) ([]*github.PullRequestReview, error)
	AddReview(ctx context.Context, id id.PR, summary, event string) error
	EnableAutoMerge(ctx context.Context, id id.PR, method githubv4.PullRequestMergeMethod) error
	IssueComment(ctx context.Context, id id.PR, comment string) error
	IssueCommentForError(ctx context.Context, id id.PR, err pe.APIError) error
	ListAllTopics(ctx context.Context, id id.PR) ([]string, error)
	ListRequiredStatusChecks(ctx context.Context, id id.PR, branch string) ([]string, error)
	ListFilesInRootDir(ctx context.Context, id id.PR, branch string) ([]string, error)
	ListFilesChangedInPR(ctx context.Context, id id.PR) ([]*github.CommitFile, error)
	GetBranchProtection(ctx context.Context, id id.PR, branch string) (*github.Protection, error)
	ListNamesOfFilesChangedInPR(ctx context.Context, id id.PR) ([]string, error)
}
