package data

import (
	"github.com/go-chi/chi/v5"
	prbot "github.com/marqeta/pr-bot"
	store "github.com/marqeta/pr-bot/datastore"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

type endpoint struct {
	controller *controller
}

type controller struct {
	metrics metrics.Emitter
	dao     store.Dao
}

func NewEndpoint(dao store.Dao, m metrics.Emitter) prbot.Endpoint {
	return &endpoint{
		controller: &controller{
			metrics: m,
			dao:     dao,
		},
	}
}

func (e *endpoint) Path() string {
	return prbot.APIVersion + "/data"
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
	r.Post("/{service}/pr/{owner}/{repo}/{number}", e.controller.HandleEvent)
	return r
}
