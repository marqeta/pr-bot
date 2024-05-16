package webhook

import (
	"net/http"

	"github.com/google/go-github/v50/github"
)

//go:generate mockery --name Parser --testonly
type Parser interface {
	ValidatePayload(r *http.Request, secretToken []byte) (payload []byte, err error)
	ParseWebHook(messageType string, payload []byte) (interface{}, error)
}

type ghEventsParser struct {
}

func NewGHEventsParser() Parser {
	return &ghEventsParser{}
}

func (p *ghEventsParser) ValidatePayload(r *http.Request, secretToken []byte) (payload []byte, err error) {
	return github.ValidatePayload(r, secretToken)
}

func (p *ghEventsParser) ParseWebHook(messageType string, payload []byte) (interface{}, error) {
	return github.ParseWebHook(messageType, payload)
}
