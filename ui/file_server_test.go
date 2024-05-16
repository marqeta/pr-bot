package ui

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
)

func TestStaticEmbed(t *testing.T) {
	f, err := content.Open("static/hello.txt")
	assert.Nil(t, err)
	defer f.Close()

	c, err := io.ReadAll(f)
	assert.Nil(t, err)

	content := strings.Trim(string(c), "\n")
	assert.Equal(t, "hello", content)
}

func TestFileServer(t *testing.T) {
	mock := evaluation.NewMockManager(t)
	endpoint := NewEndpoint(mock, metrics.NewNoopEmitter())
	assert.Equal(t, "/ui", endpoint.Path())

	req := httptest.NewRequest("GET", "http://example.com/ui/static/hello.txt", nil)
	res := executeRequest(req, endpoint)

	assert.Equal(t, http.StatusOK, res.Code)
	content := strings.Trim(res.Body.String(), "\n")
	assert.Equal(t, "hello", content)
}

func executeRequest(req *http.Request, e prbot.Endpoint) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Mount(e.Path(), e.Routes())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
