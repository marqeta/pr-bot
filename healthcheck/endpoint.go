package healthcheck

import (
	"github.com/go-chi/chi/v5"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
)

type endpoint struct {
	controller controller
}

type controller struct {
	metrics metrics.Emitter
}

func NewEndpoint(m metrics.Emitter) prbot.Endpoint {
	return &endpoint{
		controller: controller{
			metrics: m,
		},
	}
}

func (s *endpoint) Path() string {
	return prbot.APIVersion + "/health"
}

func (s *endpoint) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", s.controller.Healthcheck)

	return r
}
