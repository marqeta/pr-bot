package rate

import (
	"context"

	"github.com/marqeta/pr-bot/id"
)

type Throttler interface {
	ShouldThrottle(ctx context.Context, ID id.PR) error
	Name() string
	Key(id id.PR) string
}

type Mock struct {
	err error
}

// Key implements Throttler
func (m *Mock) Key(_ id.PR) string {
	return "mock"
}

// Name implements Throttler
func (m *Mock) Name() string {
	return "mock"
}

// ShouldThrottle implements Throttler
func (m *Mock) ShouldThrottle(_ context.Context, _ id.PR) error {
	return m.err
}

func NewMockThrottler(err error) Throttler {
	return &Mock{
		err: err,
	}
}
