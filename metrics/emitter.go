package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/slok/go-http-metrics/metrics"
)

type Emitter interface {
	EmitDist(ctx context.Context, name string, value float64, tags []string)
	EmitGauge(ctx context.Context, name string, value float64, tags []string)
	metrics.Recorder
	Close()
}

type datadogEmitter struct {
	client statsd.ClientInterface
}

func NewEmitter(client statsd.ClientInterface) Emitter {
	return &datadogEmitter{
		client: client,
	}
}

func NewNoopEmitter() Emitter {
	return &datadogEmitter{
		client: &statsd.NoOpClient{},
	}
}

// EmitGauge implements Emitter
func (d *datadogEmitter) EmitGauge(ctx context.Context, name string, value float64, tags []string) {
	oplog := httplog.LogEntry(ctx)
	err := d.client.Gauge(name, value, tags, 1)
	if err != nil {
		oplog.Err(err).Str("metric", name).Float64("metricValue", value).
			Msgf("error submitting metric: %s", name)
	}
}

// Publish implements Publisher
func (d *datadogEmitter) Close() {
	err := d.client.Close()
	if err != nil {
		log.Err(err).Msg("error closing datadog statsd client")
	}
}

// Publish implements Publisher
func (d *datadogEmitter) EmitDist(ctx context.Context, name string, value float64, tags []string) {
	oplog := httplog.LogEntry(ctx)
	err := d.client.Distribution(name, value, tags, 1)
	if err != nil {
		oplog.Err(err).Str("metric", name).Float64("metricValue", value).
			Msgf("error submitting metric: %s", name)
	}
}

func toTags(p metrics.HTTPReqProperties) []string {
	return []string{fmt.Sprintf("handler:%s", p.ID),
		fmt.Sprintf("code:%s", p.Code),
		fmt.Sprintf("method:%s", p.Method)}
}

func (d *datadogEmitter) ObserveHTTPRequestDuration(ctx context.Context, p metrics.HTTPReqProperties, duration time.Duration) {
	d.EmitDist(ctx, "responseTime", float64(duration.Milliseconds()), toTags(p))
}

func (d *datadogEmitter) ObserveHTTPResponseSize(ctx context.Context, p metrics.HTTPReqProperties, sizeBytes int64) {
	d.EmitDist(ctx, "responseSize", float64(sizeBytes), toTags(p))
}

func (d *datadogEmitter) AddInflightRequests(ctx context.Context, p metrics.HTTPProperties, quantity int) {
	d.EmitDist(ctx, "inFlightRequests", float64(quantity), []string{fmt.Sprintf("handler:%s", p.ID)})
}
