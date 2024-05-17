package webhook_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v50/github"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest"
	"github.com/marqeta/pr-bot/webhook"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

var errRandom = errors.New("random error")

func Test_controller_HandleEvent(t *testing.T) {
	type args struct {
		requestID      string
		deliveryID     string
		eventName      string
		event          any
		webhookSecrect string
	}
	type mockArgs struct {
		ValidatePayloadErr error
		ParseWebHookErr    error
		DipatchErr         error
	}
	tests := []struct {
		name     string
		args     args
		mockArgs mockArgs
		wantCode int
	}{
		{
			name: "should successfully handle pull request event",
			args: args{
				requestID:      "request1",
				deliveryID:     "delivery1",
				eventName:      pullrequest.EventName,
				event:          &github.PullRequestEvent{},
				webhookSecrect: "random secret",
			},
			wantCode: http.StatusAccepted,
		},
		{
			name: "should handle error from dispatcher",
			args: args{
				requestID:      "request1",
				deliveryID:     "",
				eventName:      pullrequest.EventName,
				event:          &github.PullRequestEvent{},
				webhookSecrect: "random secret",
			},
			mockArgs: mockArgs{
				ValidatePayloadErr: nil,
				ParseWebHookErr:    nil,
				DipatchErr:         errRandom,
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "should handle ParseWebHook error",
			args: args{
				requestID:      "request1",
				deliveryID:     "",
				eventName:      pullrequest.EventName,
				event:          &github.PullRequestEvent{},
				webhookSecrect: "random secret",
			},
			mockArgs: mockArgs{
				ValidatePayloadErr: nil,
				ParseWebHookErr:    errRandom,
				DipatchErr:         nil,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should handle validatePayload error",
			args: args{
				requestID:      "request1",
				deliveryID:     "",
				eventName:      pullrequest.EventName,
				event:          &github.PullRequestEvent{},
				webhookSecrect: "random secret",
			},
			mockArgs: mockArgs{
				ValidatePayloadErr: errRandom,
				ParseWebHookErr:    nil,
				DipatchErr:         nil,
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := webhook.NewMockParser(t)
			dispatcher := pullrequest.NewMockDispatcher(t)
			e := webhook.NewEndpoint(tt.args.webhookSecrect, parser, dispatcher, metrics.NewNoopEmitter())

			req, payload := NewRequest(t, tt.args.requestID, tt.args.deliveryID, tt.args.eventName, tt.args.event)

			parser.EXPECT().ValidatePayload(mock.AnythingOfType("*http.Request"), []byte(tt.args.webhookSecrect)).
				Return(payload, tt.mockArgs.ValidatePayloadErr).Once()
			if tt.mockArgs.ValidatePayloadErr == nil {
				parser.EXPECT().ParseWebHook(tt.args.eventName, payload).
					Return(tt.args.event, tt.mockArgs.ParseWebHookErr).Once()
			}
			if tt.mockArgs.ParseWebHookErr == nil && tt.mockArgs.ValidatePayloadErr == nil {
				dispatcher.EXPECT().Dispatch(mock.Anything, tt.args.deliveryID, tt.args.eventName, tt.args.event).
					Return(tt.mockArgs.DipatchErr).Once()
			}

			res := executeRequest(req, e)

			assert.Equal(t, tt.wantCode, res.Code)

			var wantRes webhook.EventResponse
			err := json.Unmarshal(res.Body.Bytes(), &wantRes)
			assert.Nil(t, err, "error when unmarshalling response")
			assert.Equal(t, tt.args.deliveryID, wantRes.DeliveryID)
			assert.Equal(t, tt.args.requestID, wantRes.RequestID)
		})
	}
}

func NewRequest(t *testing.T, reqID string, deliveryID string, eventName string, event any) (*http.Request, []byte) {
	body, err := json.Marshal(event)
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	r, err := http.NewRequest("POST", "http://localhost/v1/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	r.Header.Set(github.DeliveryIDHeader, deliveryID)
	r.Header.Set(github.EventTypeHeader, eventName)
	r.Header.Set(middleware.RequestIDHeader, reqID)
	return r, body
}

func executeRequest(req *http.Request, e prbot.Endpoint) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Mount(e.Path(), e.Routes())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
