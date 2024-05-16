package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
)

//go:embed static
var content embed.FS

type endpoint struct {
	controller controller
}

func NewEndpointFromFS(dir http.FileSystem, manager evaluation.Manager, m metrics.Emitter) prbot.Endpoint {
	return &endpoint{
		controller: controller{
			metrics:           m,
			staticContents:    dir,
			evaluationManager: manager,
		},
	}
}

func NewEndpoint(manager evaluation.Manager, m metrics.Emitter) prbot.Endpoint {
	files, _ := fs.Sub(content, "static")
	return NewEndpointFromFS(http.FS(files), manager, m)
}

func (s *endpoint) Path() string {
	return prbot.UI
}

func (s *endpoint) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/static/*", s.controller.FileServer)
	r.Get("/eval/{owner}/{repo}/pull/{pr}/events/{deliveryID}", s.controller.ReportDetails)
	r.Get("/eval/{owner}/{repo}/pull/{pr}", s.controller.ListReports)
	r.Get("/eval/{owner}/{repo}/{pr}/events/{deliveryID}", s.controller.ReportDetails)
	r.Get("/eval/{owner}/{repo}/{pr}", s.controller.ListReports)
	r.NotFound(s.controller.NotFound)

	return r

}
