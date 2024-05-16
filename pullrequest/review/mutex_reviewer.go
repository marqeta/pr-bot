package review

import (
	"context"
	"fmt"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/httplog"

	"github.com/marqeta/pr-bot/id"
)

//go:generate mockery --name Locker --testonly
type Locker interface {
	AcquireLockWithContext(ctx context.Context, key string, opts ...dynamolock.AcquireLockOption) (*dynamolock.Lock, error)
	ReleaseLockWithContext(ctx context.Context, lockItem *dynamolock.Lock, opts ...dynamolock.ReleaseLockOption) (bool, error)
}

func NewLocker(DDB *dynamodb.Client, table string) (Locker, error) {
	return dynamolock.New(DDB, table,
		dynamolock.WithLeaseDuration(10*time.Second),
		dynamolock.WithHeartbeatPeriod(3*time.Second))
}

type mutexReviewer struct {
	delegate Reviewer
	locker   Locker
}

func (r *mutexReviewer) releaseLock(ctx context.Context, lock *dynamolock.Lock, id id.PR) {
	oplog := httplog.LogEntry(ctx)
	success, err := r.locker.ReleaseLockWithContext(ctx, lock)
	if !success {
		oplog.Err(err).Msgf("lock for PR %v was already released", id.URL)
	}
	if err != nil {
		oplog.Err(err).Msgf("error releasing lock for PR %v", id.URL)
	}
}

func (r *mutexReviewer) acquireLock(ctx context.Context, id id.PR) (*dynamolock.Lock, error) {
	oplog := httplog.LogEntry(ctx)
	lockID := fmt.Sprintf("%v/%v", id.RepoFullName, id.Number)
	lock, err := r.locker.AcquireLockWithContext(ctx, lockID,
		dynamolock.WithRefreshPeriod(2*time.Second),
		dynamolock.WithDeleteLockOnRelease(),
	)
	if err != nil {
		oplog.Err(err).Msgf("error acquiring lock for PR %v", id.URL)
		return nil, err
	}
	return lock, nil
}

// Approve implements Reviewer.
func (r *mutexReviewer) Approve(ctx context.Context, id id.PR, body string, opts ApproveOptions) error {
	lock, err := r.acquireLock(ctx, id)
	if err != nil {
		return err
	}
	defer r.releaseLock(ctx, lock, id)
	return r.delegate.Approve(ctx, id, body, opts)
}

// Comment implements Reviewer.
func (r *mutexReviewer) Comment(ctx context.Context, id id.PR, body string) error {
	lock, err := r.acquireLock(ctx, id)
	if err != nil {
		return err
	}
	defer r.releaseLock(ctx, lock, id)
	return r.delegate.Comment(ctx, id, body)
}

// RequestChanges implements Reviewer.
func (r *mutexReviewer) RequestChanges(ctx context.Context, id id.PR, body string) error {
	lock, err := r.acquireLock(ctx, id)
	if err != nil {
		return err
	}
	defer r.releaseLock(ctx, lock, id)
	return r.delegate.RequestChanges(ctx, id, body)
}

func NewMutexReviewer(delegate Reviewer, locker Locker) Reviewer {
	return &mutexReviewer{
		delegate: delegate,
		locker:   locker,
	}
}
