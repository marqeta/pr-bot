package errors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v50/github"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/stretchr/testify/assert"
)

func TestRenderError(t *testing.T) {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "requestID")

	//nolint:goerr113
	errRandom := errors.New("random")
	apiError := pe.APIError{
		RequestID:  "requestID",
		ErrorCode:  "code",
		Message:    "msg",
		StatusCode: http.StatusBadRequest,
		ErrorText:  errRandom.Error(),
	}
	tests := []struct {
		name    string
		err     error
		wantErr pe.APIError
	}{
		{
			name: "Should render APIError",
			err:  apiError,
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "code",
				Message:    "msg",
				StatusCode: http.StatusBadRequest,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render TooManyRequestError",
			err:  pe.TooManyRequestError(ctx, "random", errRandom),
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "TooManyRequestsError",
				Message:    "random",
				StatusCode: http.StatusTooManyRequests,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render InValidRequestError",
			err:  pe.InValidRequestError(ctx, "random", errRandom),
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "InValidRequestError",
				Message:    "random",
				StatusCode: http.StatusBadRequest,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render ServiceFault",
			err:  pe.ServiceFault(ctx, "random", errRandom),
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "ServiceFault",
				Message:    "random",
				StatusCode: http.StatusInternalServerError,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render UserError",
			err:  pe.UserError(ctx, "random", errRandom),
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "UserError",
				Message:    "random",
				StatusCode: http.StatusBadRequest,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render non APIError",
			err:  errRandom,
			wantErr: pe.APIError{
				RequestID:  "requestID",
				ErrorCode:  "Unknown",
				StatusCode: http.StatusInternalServerError,
				ErrorText:  errRandom.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := NewRequest(t, tt.wantErr.RequestID)
			rr := httptest.NewRecorder()
			pe.RenderError(rr, req, tt.err)
			assert.Equal(t, tt.wantErr.StatusCode, rr.Code)
			var res pe.APIError
			err := json.Unmarshal(rr.Body.Bytes(), &res)
			assert.Nil(t, err, "error when unmarshalling response")
			assert.Equal(t, tt.wantErr, res)
		})
	}
}

func NewRequest(t *testing.T, reqID string) (*http.Request, []byte) {
	body, err := json.Marshal(&github.PullRequestEvent{})
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	r, err := http.NewRequest("POST", "http://localhost/v1/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("error creating mock request %v", err)
	}
	ctx := context.WithValue(context.TODO(), middleware.RequestIDKey, reqID)
	return r.WithContext(ctx), body
}
