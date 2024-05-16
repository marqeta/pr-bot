package rate

import (
	"context"
	"fmt"

	"github.com/go-chi/httplog"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
)

type facade struct {
	throttlers []Throttler
	metrics    metrics.Emitter
}

// Key implements Throttler
func (*facade) Key(_ id.PR) string {
	return "facade"
}

// Name implements Throttler
func (f *facade) Name() string {
	return "facade"
}

// ShouldThrottle implements Throttler
func (f *facade) ShouldThrottle(ctx context.Context, ID id.PR) error {
	oplog := httplog.LogEntry(ctx)
	for _, throttler := range f.throttlers {
		err := throttler.ShouldThrottle(ctx, ID)
		if err != nil {
			oplog.Err(err).Msgf("throttler.%v.ShouldThrottle=true for %v", throttler.Name(), throttler.Key(ID))
			t := ID.ToTags()
			t = append(t, fmt.Sprintf("throttler:%s", throttler.Name()))
			t = append(t, fmt.Sprintf("throttleKey:%s", throttler.Key(ID)))
			f.metrics.EmitDist(ctx, "throttledPRs", 1.0, t)
			return err
		}
		oplog.Info().Msgf("throttler.%v.ShouldThrottle=false for %v", throttler.Name(), throttler.Key(ID))
	}
	oplog.Info().Msgf("throttler.%v.ShouldThrottle=false", f.Name())
	return nil
}

func NewFacade(metrics metrics.Emitter, throttlers ...Throttler) Throttler {
	return &facade{
		throttlers: throttlers,
		metrics:    metrics,
	}
}
