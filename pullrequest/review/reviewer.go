package review

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-chi/httplog"
	pe "github.com/marqeta/pr-bot/errors"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/shurcooL/githubv4"
)

const (
	AutoMergeError = "Error enabling auto merge on PR"
)

type ApproveOptions struct {
	AutoMergeEnabled bool
	DefaultBranch    string
	MergeMethod      githubv4.PullRequestMergeMethod
}

//go:generate mockery --name Reviewer
type Reviewer interface {
	Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error
	RequestChanges(ctx context.Context, id id.PR, body string) error
	Comment(ctx context.Context, id id.PR, body string) error
}

type reviewer struct {
	api     gh.API
	metrics metrics.Emitter
}

// Approve implements Reviewer.
func (r *reviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	oplog := httplog.LogEntry(ctx)
	err := r.api.EnableAutoMerge(ctx, id, opts.MergeMethod)
	if err != nil {
		oplog.Err(err).Msgf("error enabling auto merge on PR %v", id.URL)
		return r.handleAutoMergeError(ctx, id, err)
	}
	oplog.Info().Msgf("enabled auto merge on PR")
	err = r.api.AddReview(ctx, id, body, gh.Approve)
	if err != nil {
		oplog.Err(err).Msgf("error approving PR")
		ae := pe.ServiceFault(ctx, "Error approving PR", err)
		// TODO publish error in UI and/or as comments on PR
		return ae
	}
	tags := append(id.ToTags(), fmt.Sprintf("mergeMethod:%s", opts.MergeMethod))
	tags = append(tags, fmt.Sprintf("reviewType:%s", "approve"))
	r.metrics.EmitDist(ctx, "reviewedPRs", 1.0, tags)
	r.metrics.EmitDist(ctx, "approvedPRs", 1.0, tags)
	oplog.Info().Msgf("reviewed PR reviewType:approve")
	return nil
}

// Comment implements Reviewer.
func (r *reviewer) Comment(ctx context.Context, id id.PR, body string) error {
	oplog := httplog.LogEntry(ctx)
	err := r.api.AddReview(ctx, id, body, gh.Comment)
	if err != nil {
		oplog.Err(err).Msgf("error reviewing PR with reviewType:comment %v", id.URL)
		ae := pe.ServiceFault(ctx, "error reviewing PR with reviewType:comment", err)
		// TODO publish error in UI and/or as comments on PR
		return ae
	}
	tags := append(id.ToTags(), fmt.Sprintf("reviewType:%s", "comment"))
	r.metrics.EmitDist(ctx, "reviewedPRs", 1.0, tags)
	r.metrics.EmitDist(ctx, "commentedPRs", 1.0, tags)
	oplog.Info().Msgf("reviewed PR reviewType:comment")
	return nil
}

// RequestChanges implements Reviewer.
func (r *reviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	oplog := httplog.LogEntry(ctx)
	err := r.api.AddReview(ctx, id, body, gh.RequestChanges)
	if err != nil {
		oplog.Err(err).Msgf("error reviewing PR with reviewType:changes_requested %v", id.URL)
		ae := pe.ServiceFault(ctx, "error reviewing PR with reviewType:changes_requested", err)
		// TODO publish error in UI and/or as comments on PR
		return ae
	}
	tags := append(id.ToTags(), fmt.Sprintf("reviewType:%s", "request_changes"))
	r.metrics.EmitDist(ctx, "reviewedPRs", 1.0, tags)
	r.metrics.EmitDist(ctx, "changesRequestedPRs", 1.0, tags)
	oplog.Info().Msgf("reviewed PR reviewType:request_changes")
	return nil
}

func NewReviewer(dao gh.API, metrics metrics.Emitter) Reviewer {
	return &reviewer{api: dao, metrics: metrics}
}

func (r *reviewer) handleAutoMergeError(ctx context.Context, id id.PR, err error) error {
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "pull request auto merge is not allowed") {
		ae := pe.UserError(ctx, AutoMergeError, err)
		r.metrics.EmitDist(ctx, "autoMergeDisabled", 1.0, id.ToTags())
		// TODO publish error in UI and/or as comments on PR
		return ae
	}
	if strings.Contains(msg, "pull request is in has_hooks status") ||
		strings.Contains(msg, "pull request is in clean status") {
		friendlyErr := fmt.Errorf("enable atleast one branch protection rule on the default branch : %w", err)
		ae := pe.UserError(ctx, AutoMergeError, friendlyErr)
		r.metrics.EmitDist(ctx, "noBranchProtectionRules", 1.0, id.ToTags())
		// TODO publish error in UI and/or as comments on PR
		return ae
	}
	ae := pe.ServiceFault(ctx, AutoMergeError, err)
	// TODO publish error in UI and/or as comments on PR
	return ae
}
