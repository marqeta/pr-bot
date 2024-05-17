package review_test

import (
	"context"
	"errors"
	"testing"

	"cirello.io/dynamolock/v2"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/stretchr/testify/mock"
)

func Test_mutex_reviewer_Approve(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(delegate *review.MockReviewer, locker *review.MockLocker)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should acquire lock, call delegate.Approve and release lock",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "random body", review.ApproveOptions{}).
						Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should acquire lock, call delegate.Approve and release lock even if approve returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "random body", review.ApproveOptions{}).
						Return(errRandom)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when locker returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *review.MockReviewer, locker *review.MockLocker) {
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should not return error when release lock returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Approve(ctx, sampleID(), "random body", review.ApproveOptions{}).
						Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(false, errRandom)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := review.NewMockReviewer(t)
			locker := review.NewMockLocker(t)
			tt.args.setExpectations(delegate, locker)
			r := review.NewMutexReviewer(delegate, locker)
			if err := r.Approve(ctx, tt.args.id, tt.args.body, review.ApproveOptions{}); (err != nil) != tt.wantErr {
				t.Errorf("mutexReviewer.Approve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mutexReviewer_Comment(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(delegate *review.MockReviewer, locker *review.MockLocker)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should acquire lock, call delegate.Comment and release lock",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "random body").Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should acquire lock, call delegate.Comment and release lock even if approve returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "random body").Return(errRandom)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when locker returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *review.MockReviewer, locker *review.MockLocker) {
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should not return error when release lock returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().Comment(ctx, sampleID(), "random body").Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(false, errRandom)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := review.NewMockReviewer(t)
			locker := review.NewMockLocker(t)
			tt.args.setExpectations(delegate, locker)
			r := review.NewMutexReviewer(delegate, locker)
			if err := r.Comment(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("rateLimitedReviewer.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mutexReviewer_RequestChanges(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")
	type args struct {
		id              id.PR
		body            string
		setExpectations func(delegate *review.MockReviewer, locker *review.MockLocker)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should acquire lock, call delegate.RequestChanges and release lock",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "random body").Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should acquire lock, call delegate.RequestChanges and release lock even if approve returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "random body").Return(errRandom)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(true, nil)
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when locker returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(_ *review.MockReviewer, locker *review.MockLocker) {
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
		{
			name: "Should not return error when release lock returns error",
			args: args{
				id:   sampleID(),
				body: "random body",
				setExpectations: func(delegate *review.MockReviewer, locker *review.MockLocker) {
					lock := &dynamolock.Lock{}
					locker.EXPECT().AcquireLockWithContext(ctx, "owner1/repo1/1",
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
						mock.AnythingOfType("dynamolock.AcquireLockOption"),
					).Return(lock, nil)
					delegate.EXPECT().RequestChanges(ctx, sampleID(), "random body").Return(nil)
					locker.EXPECT().ReleaseLockWithContext(ctx, lock).Return(false, errRandom)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := review.NewMockReviewer(t)
			locker := review.NewMockLocker(t)
			tt.args.setExpectations(delegate, locker)
			r := review.NewMutexReviewer(delegate, locker)
			if err := r.RequestChanges(ctx, tt.args.id, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("rateLimitedReviewer.RequestChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
