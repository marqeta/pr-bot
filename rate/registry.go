package rate

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jonboulle/clockwork"
	lim "github.com/mennanov/limiters"
	"github.com/rs/zerolog/log"
)

var ErrCastingIntoSWLimiter = errors.New("error type casting sliding window limiter")

//go:generate mockery --name Limiter --testonly
type Limiter interface {
	Limit(ctx context.Context) (time.Duration, error)
}

//go:generate mockery --name Getter --testonly
type Getter interface {
	Name() string
	GetOrCreate(ctx context.Context, key string, limit Limit) (Limiter, error)
	Close()
}
type swRegistry struct {
	name   string
	reg    *lim.Registry
	clock  clockwork.Clock
	ticker clockwork.Ticker
	client *dynamodb.Client
	props  lim.DynamoDBTableProperties
	done   chan bool
}

// Name implements Getter
func (sw *swRegistry) Name() string {
	return sw.name
}

// Close implements Getter
func (sw *swRegistry) Close() {
	sw.done <- true
}

// GetOrCreate implements Getter
func (sw *swRegistry) GetOrCreate(_ context.Context, key string, limit Limit) (Limiter, error) {

	l := sw.reg.GetOrCreate(key,
		func() interface{} {
			backend := lim.NewSlidingWindowDynamoDB(sw.client, key, sw.props)
			return lim.NewSlidingWindow(limit.Value, limit.Window, backend, sw.clock, 1e-4)
		},
		2*limit.Window, sw.clock.Now())

	limiter, ok := l.(*lim.SlidingWindow)
	if !ok {
		return nil, ErrCastingIntoSWLimiter
	}
	return limiter, nil
}

func NewSWRegistry(name string, client *dynamodb.Client, props lim.DynamoDBTableProperties,
	clock clockwork.Clock) Getter {
	reg := &swRegistry{
		name:   name,
		reg:    lim.NewRegistry(),
		clock:  clock,
		ticker: clock.NewTicker(5 * time.Minute),
		client: client,
		props:  props,
		done:   make(chan bool),
	}
	go reg.removeExpired()
	return reg
}

func (sw *swRegistry) removeExpired() {
	for {
		select {
		case <-sw.done:
			sw.ticker.Stop()
			return
		case <-sw.ticker.Chan():
			count := sw.reg.DeleteExpired(sw.clock.Now())
			log.Info().Msgf("deleted %v expired limiters from %v registry", count, sw.name)
		}
	}
}
