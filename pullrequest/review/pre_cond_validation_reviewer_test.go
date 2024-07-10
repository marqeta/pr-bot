package review_test

import (
	"context"
	"errors"
	"testing"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/pullrequest/review"
)

func Test_preCondValidationReviewer_Approve(t *testing.T) {
	ctx := context.Background()
	type args struct {
		id              id.PR
		body            string
		opts            review.ApproveOptions
		setExpectations func(delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate if automerge is enabled",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(delegate *review.MockReviewer) {
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should not call delegate if automerge is disabled",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch: "main",
				},
				setExpectations: func(_ *review.MockReviewer) {
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(delegate)
			r := review.NewPreCondValidationReviewer(delegate)
			if err := r.Approve(ctx, tt.args.id, tt.args.body, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_preCondValidationReviewer_Comment(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate",
			args: args{
				id:   sampleID(),
				body: "comment",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error from delegate",
			args: args{
				id:   sampleID(),
				body: "comment",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().Comment(ctx, sampleID(), "comment").Return(errRandom)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewPreCondValidationReviewer(delegate)
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_preCondValidationReviewer_RequestChanges(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate",
			args: args{
				id:   sampleID(),
				body: "changes",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "changes").Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error from delegate",
			args: args{
				id:   sampleID(),
				body: "changes",
				setExpectations: func(_ *gh.MockAPI, delegate *review.MockReviewer) {
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "changes").Return(errRandom)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewPreCondValidationReviewer(delegate)
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
