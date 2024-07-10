package review_test

import (
	"context"
	"errors"
	"testing"

	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest/review"
)

func Test_BBReviewer_Approve(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		opts            review.ApproveOptions
		setExpectations func(api *gh.MockAPI, delegate *review.MockReviewer)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should call delegate for non bb repos",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "file2"}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw error from delegate for non bb repos",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "file2"}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should call delegate for bb repos with bb-ci",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "blackbird.yaml"}, nil)
					api.EXPECT().ListRequiredStatusChecks(ctx, sampleID(), "main").
						Return([]string{"status1", "blackbird-ci"}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should throw from delegate for bb repos with bb-ci",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "blackbird.yaml"}, nil)
					api.EXPECT().ListRequiredStatusChecks(ctx, sampleID(), "main").
						Return([]string{"status1", "blackbird-ci"}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should not call delegate for bb repos without bb-ci",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListRequiredStatusChecks(ctx, sampleID(), "main").
						Return([]string{"status1", "status2"}, nil)
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "blackbird.yaml"}, nil)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error if listRequiredStatusChecks fails",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{"file1", "blackbird.yaml"}, nil)
					api.EXPECT().ListRequiredStatusChecks(ctx, sampleID(), "main").
						Return([]string{}, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should throw error if listfilesinrootdir fails",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, _ *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{}, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should call delegate on empty/new repo",
			args: args{
				id:   sampleID(),
				body: "LGTM",
				opts: review.ApproveOptions{
					DefaultBranch:    "main",
					AutoMergeEnabled: true,
				},
				setExpectations: func(api *gh.MockAPI, delegate *review.MockReviewer) {
					api.EXPECT().ListFilesInRootDir(ctx, sampleID(), "main").
						Return([]string{}, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "LGTM", review.ApproveOptions{
						DefaultBranch:    "main",
						AutoMergeEnabled: true,
					}).Return(nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			delegate := review.NewMockReviewer(t)
			tt.args.setExpectations(api, delegate)
			r := review.NewBBReviewer(delegate, api, metrics.NewNoopEmitter())
			if err := r.Approve(ctx, tt.args.id, tt.args.body, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_BBReviewer_Comment(t *testing.T) {
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
			r := review.NewBBReviewer(delegate, api, metrics.NewNoopEmitter())
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_BBReviewer_RequestChanges(t *testing.T) {
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
			r := review.NewBBReviewer(delegate, api, metrics.NewNoopEmitter())
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("preCondValidationReviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
