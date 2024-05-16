package prbot

import (
	"github.com/go-chi/chi/v5"
)

type Endpoint interface {
	Routes() chi.Router
	Path() string
}

const APIVersion string = "/v1"
const UI string = "/ui"
