package datastore

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/marqeta/pr-bot/id"
)

var (
	ErrMissingOwner    = errors.New("missing owner")
	ErrInvalidPRNumber = errors.New("missing or invalid pr number")
	ErrMissingHead     = errors.New("missing invalid head sha")
	ErrMissingBase     = errors.New("missing base sha")
	ErrMissingRepo     = errors.New("missing repo")
	ErrInvalidService  = errors.New("invalid service")
	ErrMissingJob      = errors.New("missing job")
	ErrInvalidHeadBase = errors.New("head and base sha are the same")
)

type Metadata struct {
	id.PR
	Service string `json:"service"`
	Head    string `json:"head"`
	Base    string `json:"base"`
	Job     string `json:"job,omitempty"`
}

func ToMetadata(r *http.Request) (*Metadata, error) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	numberStr := chi.URLParam(r, "number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, ErrInvalidPRNumber
	}
	service := chi.URLParam(r, "service")
	head := r.URL.Query().Get("head")
	base := r.URL.Query().Get("base")
	job := r.URL.Query().Get("job")
	fmt.Println(service)

	if service == "" {
		return nil, ErrInvalidService
	}
	if owner == "" {
		return nil, ErrMissingOwner
	}
	if repo == "" {
		return nil, ErrMissingRepo
	}
	if number <= 0 {
		return nil, ErrInvalidPRNumber
	}
	if head == "" {
		return nil, ErrMissingHead
	}
	if base == "" {
		return nil, ErrMissingBase
	}
	if head == base {
		return nil, ErrInvalidHeadBase
	}
	if job == "" {
		return nil, ErrMissingJob
	}
	return &Metadata{
		PR:      id.PR{Owner: owner, Repo: repo, Number: number},
		Service: service,
		Head:    head,
		Base:    base,
		Job:     job,
	}, nil
}
