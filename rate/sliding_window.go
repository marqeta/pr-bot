package rate

import (
	"context"
	"errors"
	"fmt"

	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/configstore"
	"github.com/marqeta/pr-bot/id"
	lim "github.com/mennanov/limiters"
)

type swLimiter struct {
	keyer    Keyer
	cfgStore configstore.Getter[*LimiterConfig]
	registry Getter
}

// Key implements Throttler
func (sw *swLimiter) Key(id id.PR) string {
	return sw.keyer(id)
}

// Name implements Throttler
func (sw *swLimiter) Name() string {
	return sw.registry.Name()
}

// ShouldThrottle implements Throttler
func (sw *swLimiter) ShouldThrottle(ctx context.Context, id id.PR) error {

	key := sw.keyer(id)

	limit, err := sw.getLimit(key)
	if err != nil {
		return err
	}

	limiter, err := sw.registry.GetOrCreate(ctx, key, limit)
	if err != nil {
		return err
	}

	waitTime, err := limiter.Limit(ctx)
	if err != nil && errors.Is(err, lim.ErrLimitExhausted) {
		// throttled
		msg := fmt.Sprintf("%v throttled request for key %v, try again in %v", sw.Name(), sw.Key(id), waitTime)
		return prbot.TooManyRequestError(ctx, msg, err)
	}
	return err
}

func (sw *swLimiter) getLimit(key string) (Limit, error) {
	cfg, err := sw.cfgStore.Get()
	if err != nil {
		return Limit{}, err
	}
	limit := cfg.Get(key)
	return limit, nil
}

func NewSlidingWindowLimiter(keyer Keyer, registry Getter,
	store configstore.Getter[*LimiterConfig]) Throttler {
	return &swLimiter{
		keyer:    keyer,
		cfgStore: store,
		registry: registry,
	}
}
