package data_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/data"
	store "github.com/marqeta/pr-bot/datastore"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_controller_HandleEvent(t *testing.T) {
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	payload := []byte(`{"key": "value"}`)
	userErr := pe.UserError(context.TODO(), "error storing payload", randomErr)
	serviceErr := pe.ServiceFault(context.TODO(), "internal server error", randomErr)

	tests := []struct {
		name            string
		setExpectations func(dao *store.MockDao, req *http.Request)
		requestID       string
		wantCode        int
		wantMessage     string
	}{
		{
			name: "should return 200 when storing payload",
			setExpectations: func(dao *store.MockDao, _ *http.Request) {
				dao.EXPECT().ToMetadata(mock.Anything, mock.AnythingOfType("*http.Request")).
					Return(randomMetadata(), nil).Once()
				dao.EXPECT().StorePayload(mock.Anything, randomMetadata(), json.RawMessage(payload)).
					Return(nil).Once()
			},
			requestID:   "123A",
			wantCode:    http.StatusOK,
			wantMessage: "payload stored successfully",
		},
		{
			name: "should return 400 when error parsing metadata",
			setExpectations: func(dao *store.MockDao, _ *http.Request) {
				dao.EXPECT().ToMetadata(mock.Anything, mock.AnythingOfType("*http.Request")).
					Return(nil, randomErr).Once()
			},
			requestID:   "123B",
			wantCode:    http.StatusBadRequest,
			wantMessage: "error parsing metadata",
		},
		{
			name: "should return 400 when user error occrs while storing payload",
			setExpectations: func(dao *store.MockDao, _ *http.Request) {
				err := userErr
				err.RequestID = "123C"
				dao.EXPECT().ToMetadata(mock.Anything, mock.AnythingOfType("*http.Request")).
					Return(randomMetadata(), nil).Once()
				dao.EXPECT().StorePayload(mock.Anything, randomMetadata(), json.RawMessage(payload)).
					Return(err).Once()
			},
			requestID:   "123C",
			wantCode:    http.StatusBadRequest,
			wantMessage: "error storing payload",
		},
		{
			name: "should return 500 when error occrs while storing payload",
			setExpectations: func(dao *store.MockDao, _ *http.Request) {
				err := serviceErr
				err.RequestID = "123D"
				dao.EXPECT().ToMetadata(mock.Anything, mock.AnythingOfType("*http.Request")).
					Return(randomMetadata(), nil).Once()
				dao.EXPECT().StorePayload(mock.Anything, randomMetadata(), json.RawMessage(payload)).
					Return(err).Once()
			},
			requestID:   "123D",
			wantCode:    http.StatusInternalServerError,
			wantMessage: "internal server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := store.NewMockDao(t)
			e := data.NewEndpoint(dao, metrics.NewNoopEmitter())

			req := NewRequest(t, tt.requestID, randomMetadata(), payload)
			tt.setExpectations(dao, req)
			res := executeRequest(req, e)
			assert.Equal(t, tt.wantCode, res.Code)

			var wantRes data.Response
			err := json.Unmarshal(res.Body.Bytes(), &wantRes)
			assert.Nil(t, err, "error when unmarshalling response")
			assert.Equal(t, tt.requestID, wantRes.RequestID)
			assert.Equal(t, tt.wantMessage, wantRes.Message)
			assert.Equal(t, tt.wantCode, wantRes.StatusCode)
		})
	}
}

func randomMetadata() *store.Metadata {
	return &store.Metadata{
		PR: id.PR{
			Owner:  "ci",
			Repo:   "ci-kirkland-jobs",
			Number: 1,
		},
		Service: "kirkland",
		Head:    "asdahasdsdjasdjacn1123h12",
		Base:    "1233511adsfaewdqffasdasda",
		Job:     "ci-kirkland-jobs (us-west-2))",
	}
}

func executeRequest(req *http.Request, e prbot.Endpoint) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Mount(e.Path(), e.Routes())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func NewRequest(t *testing.T, reqID string, m *store.Metadata, payload []byte) *http.Request {
	url := fmt.Sprintf("http://localhost/v1/data/%s/pr/%s/%s/%d?head=%s&base=%s&job=%s",
		m.Service, m.Owner, m.Repo, m.Number, m.Head, m.Base, m.Job)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	r.Header.Set(middleware.RequestIDHeader, reqID)
	return r
}
