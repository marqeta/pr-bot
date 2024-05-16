package configstore

import (
	"fmt"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/marqeta/pr-bot/metrics"

	"github.com/jonboulle/clockwork"
)

type dbStore[T DynamicConfig] struct {
	dao       Dao[T]
	ticker    clockwork.Ticker
	config    atomic.Value
	done      chan bool
	tableName string
	metrics   metrics.Emitter
	name      string
	zero      T
}

// Get implements Getter
func (d *dbStore[T]) Get() (T, error) {
	v, ok := d.config.Load().(T)
	if !ok {
		//nolint:goerr113
		return d.zero, fmt.Errorf("error type casting value: %#v into type %s ", v, d.name)
	}
	return v, nil
}

// Close implements Getter
func (d *dbStore[T]) Close() {
	d.done <- true
}

func NewDBStore[T DynamicConfig](dao Dao[T], name string, tableName string, ticker clockwork.Ticker, metrics metrics.Emitter) (Getter[T], error) {
	db := &dbStore[T]{
		dao:       dao,
		tableName: tableName,
		ticker:    ticker,
		metrics:   metrics,
		done:      make(chan bool),
		name:      name,
	}
	err := db.loadOnce()
	if err != nil {
		log.Err(err).Msgf("error loading dynamic config of type %s", name)
		return nil, err
	}
	go db.load()
	return db, nil
}

func (d *dbStore[T]) loadOnce() error {
	cfg, err := d.dao.GetItem(d.name, d.tableName)
	if err != nil {
		return err
	}
	err = cfg.Update()
	if err != nil {
		emitError(d.metrics, d.name, "UpdateHookError")
		return err
	}
	d.config.Store(cfg)
	emitSuccess(d.metrics, d.name)
	log.Info().Interface("Config", cfg).Msgf("successfully loaded dynamic config of type %s", d.name)
	return nil
}

func (d *dbStore[T]) load() {
	for {
		select {
		case <-d.done:
			d.ticker.Stop()
			return
		case <-d.ticker.Chan():
			err := d.loadOnce()
			if err != nil {
				log.Err(err).Msgf("error loading dynamic config of type %s", d.name)
			}
		}
	}
}
