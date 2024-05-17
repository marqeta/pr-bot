package webhook

import (
	"github.com/go-chi/chi/v5"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

type endpoint struct {
	controller *controller
}

type controller struct {
	webhookSecrect string
	parser         Parser
	dispatcher     pullrequest.Dispatcher
	metrics        metrics.Emitter
}

func NewEndpoint(ws string, p Parser, d pullrequest.Dispatcher, m metrics.Emitter) prbot.Endpoint {
	return &endpoint{
		controller: &controller{
			webhookSecrect: ws,
			parser:         p,
			dispatcher:     d,
			metrics:        m,
		},
	}
}

func (e *endpoint) Path() string {
	return prbot.APIVersion + "/webhook"
}

func (e *endpoint) Routes() chi.Router {
	r := chi.NewRouter()
	// publishes handler metrics i.e. responseTime
	mdlw := middleware.New(middleware.Config{
		Recorder:               e.controller.metrics,
		GroupedStatus:          true,
		DisableMeasureSize:     true,
		DisableMeasureInflight: true,
	})
	// use handler url as the id
	r.Use(std.HandlerProvider("", mdlw))
	r.Post("/", e.controller.HandleEvent)
	return r
}
