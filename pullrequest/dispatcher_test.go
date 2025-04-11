package pullrequest_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest"
)

var errRandom = errors.New("random error")

func Test_dispatcher_Dispatch(t *testing.T) {
	ctx := context.Background()
	type args struct {
		deliveryID string
		eventName  string
		event      *github.PullRequestEvent
	}
	tests := []struct {
		name            string
		args            args
		setExpectations func(id id.PR, event *github.PullRequestEvent,
			f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler)
		wantErr error
	}{
		{
			name: "Should dispatch PR opened Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("opened"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR reopened Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("reopened"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR edited Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("edited"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR labeled Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      labeled(prEvent(github.String("labeled"), sampleID())),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR unlabeled Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("unlabeled"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR review requested Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("review_requested"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR review_request_removed Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("review_request_removed"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR assigned Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("assigned"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR unassigned Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("unassigned"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR synchronize Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("synchronize"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should dispatch PR edited Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("edited"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should return silently when PR event action is unknown",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("unknown"), sampleID()),
			},
			setExpectations: func(id id.PR, _ *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should skip dispatch when shouldHandle=false",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("unknown"), sampleID()),
			},
			setExpectations: func(id id.PR, _ *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(false, nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Error when EventName != pull_request",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName + "asd asd",
				event:      prEvent(github.String("reopened"), sampleID()),
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: pe.InValidRequestError(ctx, "error parsing webhook event", pullrequest.ErrMismatchedEvent),
		},
		{
			name: "Error when Event action is nil",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(nil, sampleID()),
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: pe.InValidRequestError(ctx, "error parsing webhook event", pullrequest.ErrEventActionNotFound),
		},
		{
			name: "Error when Event PR is nil",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event: &github.PullRequestEvent{
					Action:      github.String("reopened"),
					PullRequest: nil,
				},
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: pe.InValidRequestError(ctx, "error parsing webhook event", pullrequest.ErrPRNotFound),
		},
		{
			name: "Should return error from event handler",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("reopened"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(errRandom).Once()
			},
			wantErr: errRandom,
		},
		{
			name: "Should return error from event handler opened event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("opened"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(errRandom).Once()
			},
			wantErr: errRandom,
		},
		{
			name: "Should return error from event handler labeled event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("labeled"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPREvent(ctx, id, event).Return(errRandom).Once()
			},
			wantErr: errRandom,
		},
		{
			name: "Should return error from shouldHandle",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("reopened"), sampleID()),
			},
			setExpectations: func(id id.PR, _ *github.PullRequestEvent,
				f *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(false, errRandom).Once()
			},
			wantErr: errRandom,
		},
		{
			name: "Skip dispatch for private repo PR event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      private(prEvent(github.String("labeled"), sampleID())),
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: nil,
		},
		{
			name: "Skip dispatch when event.Repo.Visibility is nil",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("labeled"), sampleID()),
			},
			setExpectations: func(_ id.PR, event *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				event.Repo.Visibility = nil
			},
			wantErr: nil,
		},
		{
			name: "Skip dispatch when event.Repo.Visibility is empty",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventName,
				event:      prEvent(github.String("labeled"), sampleID()),
			},
			setExpectations: func(_ id.PR, event *github.PullRequestEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				event.Repo.Visibility = github.String("")
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := pullrequest.NewMockEventHandler(t)
			filter := pullrequest.NewMockEventFilter(t)
			d := pullrequest.NewDispatcher(handler, filter, metrics.NewNoopEmitter())

			tt.setExpectations(sampleID(), tt.args.event, filter, handler)
			err := d.Dispatch(ctx, tt.args.deliveryID, tt.args.eventName, tt.args.event)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("dispatcher.Dispatch() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func sampleID() id.PR {

	return id.PR{
		Owner:        "owner1",
		Repo:         "repo1",
		Number:       1,
		NodeID:       "nodeid1",
		RepoFullName: "owner1/repo1",
		Author:       "user1",
		URL:          "owner1/repo1/1",
	}
}

func Test_dispatcher_DispatchReview(t *testing.T) {
	ctx := context.Background()
	type args struct {
		deliveryID string
		eventName  string
		event      *github.PullRequestReviewEvent
	}
	tests := []struct {
		name            string
		args            args
		setExpectations func(id id.PR, event *github.PullRequestReviewEvent,
			f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler)
		wantErr error
	}{
		{
			name: "Should dispatch PR Review Submitted Event",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventNameReview,
				event:      prReviewEvent(github.String("submitted"), sampleID()),
			},
			setExpectations: func(id id.PR, event *github.PullRequestReviewEvent,
				f *pullrequest.MockEventFilter, h *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
				h.EXPECT().EvalAndReviewPRReviewEvent(ctx, id, event).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Should return silently when PR event action is unknown",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventNameReview,
				event:      prReviewEvent(github.String("unknown"), sampleID()),
			},
			setExpectations: func(id id.PR, _ *github.PullRequestReviewEvent,
				f *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
				f.EXPECT().ShouldHandle(ctx, id).Return(true, nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Error when Event action is nil",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventNameReview,
				event:      prReviewEvent(nil, sampleID()),
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestReviewEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: pe.InValidRequestError(ctx, "error parsing webhook event", pullrequest.ErrEventActionNotFound),
		},
		{
			name: "Error when Event PR is nil",
			args: args{
				deliveryID: "123",
				eventName:  pullrequest.EventNameReview,
				event: &github.PullRequestReviewEvent{
					Action:      github.String("submitted"),
					PullRequest: nil,
				},
			},
			setExpectations: func(_ id.PR, _ *github.PullRequestReviewEvent,
				_ *pullrequest.MockEventFilter, _ *pullrequest.MockEventHandler) {
			},
			wantErr: pe.InValidRequestError(ctx, "error parsing webhook event", pullrequest.ErrPRNotFound),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := pullrequest.NewMockEventHandler(t)
			filter := pullrequest.NewMockEventFilter(t)
			d := pullrequest.NewDispatcher(handler, filter, metrics.NewNoopEmitter())

			tt.setExpectations(sampleID(), tt.args.event, filter, handler)
			err := d.DispatchReview(ctx, tt.args.deliveryID, tt.args.eventName, tt.args.event)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("dispatcher.Dispatch() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func empty(e *github.PullRequestEvent) *github.PullRequestEvent {
	e.PullRequest.ChangedFiles = github.Int(0)
	return e
}
func prReviewEvent(action *string, id id.PR) *github.PullRequestReviewEvent {
	url := fmt.Sprintf("%s/%s/%d", id.Owner, id.Repo, id.Number)
	return &github.PullRequestReviewEvent{
		Action: action,
		Review: &github.PullRequestReview{
			State: github.String("APPROVED"),
			User: &github.User{
				Login: &id.Author,
			},
		},
		PullRequest: &github.PullRequest{
			Number: &id.Number,
			NodeID: &id.NodeID,
			User: &github.User{
				Login: &id.Author,
			},
			HTMLURL: &url,
			Base: &github.PullRequestBranch{
				Repo: &github.Repository{
					AllowAutoMerge:   aws.Bool(false),
					AllowRebaseMerge: aws.Bool(false),
					AllowSquashMerge: aws.Bool(false),
					AllowMergeCommit: aws.Bool(false),
				},
			},
			ChangedFiles: github.Int(1),
		},
		Repo: &github.Repository{
			Owner: &github.User{
				Login: &id.Owner,
			},
			DefaultBranch: github.String("main"),
			Name:          &id.Repo,
			FullName:      &id.RepoFullName,
			Visibility:    github.String("public"),
		},
	}
}

func prEvent(action *string, id id.PR) *github.PullRequestEvent {
	url := fmt.Sprintf("%s/%s/%d", id.Owner, id.Repo, id.Number)
	return &github.PullRequestEvent{
		Action: action,
		PullRequest: &github.PullRequest{
			Number: &id.Number,
			NodeID: &id.NodeID,
			User: &github.User{
				Login: &id.Author,
			},
			HTMLURL: &url,
			Base: &github.PullRequestBranch{
				Repo: &github.Repository{
					AllowAutoMerge:   aws.Bool(false),
					AllowRebaseMerge: aws.Bool(false),
					AllowSquashMerge: aws.Bool(false),
					AllowMergeCommit: aws.Bool(false),
				},
			},
			ChangedFiles: github.Int(1),
		},
		Repo: &github.Repository{
			Owner: &github.User{
				Login: &id.Owner,
			},
			DefaultBranch: github.String("main"),
			Name:          &id.Repo,
			FullName:      &id.RepoFullName,
			Visibility:    github.String("public"),
		},
	}
}

func allMergeMethods(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.PullRequest.Base.Repo.AllowAutoMerge = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowMergeCommit = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowRebaseMerge = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowSquashMerge = aws.Bool(true)
	return event
}

func merge(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.PullRequest.Base.Repo.AllowAutoMerge = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowMergeCommit = aws.Bool(true)
	return event
}

func rebase(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.PullRequest.Base.Repo.AllowAutoMerge = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowRebaseMerge = aws.Bool(true)
	return event
}

func squash(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.PullRequest.Base.Repo.AllowAutoMerge = aws.Bool(true)
	event.PullRequest.Base.Repo.AllowSquashMerge = aws.Bool(true)
	return event
}

func labeled(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.Label = &github.Label{
		Name: github.String("l1"),
	}
	return event
}

func private(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.Repo.Visibility = github.String("private")
	return event
}

func randomBranchAsDefault(event *github.PullRequestEvent) *github.PullRequestEvent {
	event.Repo.DefaultBranch = github.String("random")
	return event
}
