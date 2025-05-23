package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v50/github"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/shurcooL/githubv4"
)

type ApprovalMessage struct {
	RequestID string `json:"request_id"`
}

type ErrorMessage struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

type githubDao struct {
	v3         *github.Client
	v4         *githubv4.Client
	metrics    metrics.Emitter
	serverHost string
	serverPort int
}

func NewAPI(serverHost string, serverPort int, v3 *github.Client, v4 *githubv4.Client, metrics metrics.Emitter) API {
	return &githubDao{
		v3:         v3,
		v4:         v4,
		metrics:    metrics,
		serverHost: serverHost,
		serverPort: serverPort,
	}
}

// ListReviews implements Dao
func (gh *githubDao) ListReviews(ctx context.Context, id id.PR) ([]*github.PullRequestReview, error) {
	reviews, resp, err := gh.v3.PullRequests.ListReviews(ctx,
		id.Owner, id.Repo, id.Number, &github.ListOptions{
			Page:    1,
			PerPage: 50,
		})
	if err != nil {
		return nil, err
	}
	gh.emitTokenExpiration(ctx, resp)
	return reviews, nil
}

// GetBranchProtection implements Dao
func (gh *githubDao) GetBranchProtection(ctx context.Context, id id.PR, branch string) (*github.Protection, error) {
	b, resp, err := gh.v3.Repositories.GetBranchProtection(ctx, id.Owner, id.Repo, branch)
	if err != nil {
		return nil, classifyError(ctx, resp, fmt.Sprintf("error getting branch protection for PR %v", id.URL), err)
	}
	gh.emitTokenExpiration(ctx, resp)
	return b, nil
}

// GetPR implements Dao.
// Only returns first 30 file changes in PR
// TODO: iterate over pages
func (gh *githubDao) ListFilesChangedInPR(ctx context.Context, id id.PR) ([]*github.CommitFile, error) {
	files, resp, err := gh.v3.PullRequests.ListFiles(ctx, id.Owner, id.Repo, id.Number, nil)
	if err != nil {
		return nil, err
	}
	gh.emitTokenExpiration(ctx, resp)
	return files, nil
}

// ListFilesInRootDir implements Dao.
func (gh *githubDao) ListFilesInRootDir(ctx context.Context, id id.PR, branch string) ([]string, error) {
	// empty path == root dir
	_, files, resp, err := gh.v3.Repositories.GetContents(ctx, id.Owner, id.Repo, "", &github.RepositoryContentGetOptions{
		Ref: branch,
	})
	if err != nil {
		return nil, classifyError(ctx, resp, fmt.Sprintf("error listing files on repo for PR %v", id.URL), err)
	}
	gh.emitTokenExpiration(ctx, resp)
	filenames := make([]string, len(files))
	for _, f := range files {
		if f.Type != nil && *f.Type == "file" && f.Path != nil {
			filename := *f.Path
			filenames = append(filenames, filename)
		}
	}
	return filenames, nil
}

// ListRequiredStatusChecks implements Dao.
func (gh *githubDao) ListRequiredStatusChecks(ctx context.Context, id id.PR, branch string) ([]string, error) {
	checks, resp, err := gh.v3.Repositories.ListRequiredStatusChecksContexts(ctx, id.Owner, id.Repo, branch)
	if err != nil {
		return nil, classifyError(ctx, resp,
			fmt.Sprintf("error listing required status checks on repo for PR %v", id.URL), err)
	}
	gh.emitTokenExpiration(ctx, resp)
	return checks, nil
}

// ListAllTopics implements Dao
func (gh *githubDao) ListAllTopics(ctx context.Context, id id.PR) ([]string, error) {
	topics, resp, err := gh.v3.Repositories.ListAllTopics(ctx, id.Owner, id.Repo)
	if err != nil {
		return nil, err
	}
	gh.emitTokenExpiration(ctx, resp)
	return topics, err
}

// ListNamesOfFilesChangedInPR implements Dao
// TODO iterate over pages
func (gh *githubDao) ListNamesOfFilesChangedInPR(ctx context.Context, id id.PR) ([]string, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				Files struct {
					Edges []struct {
						Node struct {
							Path githubv4.String
						}
					}
				} `graphql:"files(first: 100)"`
			} `graphql:"pullRequest(number: $prNumber)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(id.Owner),
		"repo":  githubv4.String(id.Repo),
		// TODO handle int32 vs int usage in v3 and v4 version of github lcient
		//nolint:gosec
		"prNumber": githubv4.Int(id.Number),
	}
	err := gh.v4.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}
	filepaths := make([]string, 0, len(q.Repository.PullRequest.Files.Edges))
	for _, edge := range q.Repository.PullRequest.Files.Edges {
		filepaths = append(filepaths, string(edge.Node.Path))
	}
	return filepaths, nil
}

// EnableAutoMerge implements Dao
func (gh *githubDao) EnableAutoMerge(ctx context.Context, id id.PR, method githubv4.PullRequestMergeMethod) error {
	var mutation struct {
		EnablePullRequestAutoMerge struct {
			PullRequest struct {
				Title githubv4.String
			}
		} `graphql:"enablePullRequestAutoMerge(input: $input)"`
	}
	input := githubv4.EnablePullRequestAutoMergeInput{
		PullRequestID: id.NodeID,
		MergeMethod:   &method,
	}
	err := gh.v4.Mutate(ctx, &mutation, input, nil)
	if err != nil {
		return err
	}
	return nil
}

// IssueComment implements Dao
func (gh *githubDao) IssueComment(ctx context.Context, id id.PR, comment string) error {
	_, resp, err := gh.v3.Issues.CreateComment(ctx, id.Owner, id.Repo, id.Number,
		&github.IssueComment{
			Body: &comment,
		})
	if err != nil {
		return err
	}
	gh.emitTokenExpiration(ctx, resp)
	return nil
}

// AddReview implements Dao
func (gh *githubDao) AddReview(ctx context.Context, id id.PR, summary, event string) error {
	msg := ApprovalMessage{
		RequestID: middleware.GetReqID(ctx),
	}
	b, e := json.MarshalIndent(msg, "", "  ")
	if e != nil {
		return e
	}
	body := fmt.Sprintf(ApprovalTemplate, summary, gh.UI(id), string(b))
	_, resp, err := gh.v3.PullRequests.CreateReview(ctx, id.Owner, id.Repo, id.Number,
		&github.PullRequestReviewRequest{
			Body:  &body,
			Event: &event,
		})
	if err != nil {
		return err
	}
	gh.emitTokenExpiration(ctx, resp)
	return nil
}

// IssueCommentForError implements Dao
func (gh *githubDao) IssueCommentForError(ctx context.Context, id id.PR, apiError pe.APIError) error {
	b, err := json.MarshalIndent(apiError, "", "  ")
	if err != nil {
		return err
	}
	err = gh.IssueComment(ctx, id, fmt.Sprintf(ErrorTemplate, apiError.Message, string(b)))
	if err != nil {
		return err
	}
	return nil
}

func (gh *githubDao) GetPullRequest(ctx context.Context, id id.PR) (*github.PullRequest, error) {
	pr, resp, err := gh.v3.PullRequests.Get(ctx, id.Owner, id.Repo, id.Number)
	if err != nil {
		return nil, classifyError(ctx, resp, "error getting PR details", err)
	}
	gh.emitTokenExpiration(ctx, resp)
	return pr, nil
}

func (gh *githubDao) GetRepository(ctx context.Context, id id.PR) (*github.Repository, error) {
	repo, resp, err := gh.v3.Repositories.Get(ctx, id.Owner, id.Repo)
	if err != nil {
		return nil, classifyError(ctx, resp, "error getting repo details", err)
	}
	gh.emitTokenExpiration(ctx, resp)
	return repo, nil
}

func (gh *githubDao) GetOrganization(ctx context.Context, id id.PR) (*github.Organization, error) {
	org, resp, err := gh.v3.Organizations.Get(ctx, id.Owner)
	if err != nil {
		return nil, classifyError(ctx, resp, "error getting org details", err)
	}
	gh.emitTokenExpiration(ctx, resp)
	return org, nil
}

func (gh *githubDao) emitTokenExpiration(ctx context.Context, resp *github.Response) {
	if resp == nil {
		return
	}
	d := time.Until(resp.TokenExpiration.Time)
	days := d.Hours() / 24
	gh.metrics.EmitGauge(ctx, "GHETokenExpiry", days, nil)
}

func (gh *githubDao) UI(id id.PR) string {
	if gh.serverHost == "localhost" {
		return fmt.Sprintf("http://%s:%d/ui/eval/%s/pull/%d", gh.serverHost, gh.serverPort, id.RepoFullName, id.Number)
	}
	return fmt.Sprintf("https://%s/ui/eval/%s/pull/%d", gh.serverHost, id.RepoFullName, id.Number)
}

func classifyError(ctx context.Context, resp *github.Response, msg string, err error) error {
	if resp == nil {
		return err
	}
	if isClientError(resp) {
		return pe.UserError(ctx, msg, err)
	}
	return pe.ServiceFault(ctx, msg, err)
}

func isClientError(resp *github.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 400 && resp.StatusCode < 500
}
