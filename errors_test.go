package prbot_test

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
	"github.com/stretchr/testify/assert"
	prbot "github.com/marqeta/pr-bot"
)

func TestRenderError(t *testing.T) {

	//nolint:goerr113
	errRandom := errors.New("random")
	apiError := prbot.APIError{
		RequestID:  "requestID",
		ErrorCode:  "code",
		Message:    "msg",
		StatusCode: http.StatusBadRequest,
		ErrorText:  errRandom.Error(),
	}
	tests := []struct {
		name    string
		err     error
		wantErr prbot.APIError
	}{
		{
			name: "Should render APIError",
			err:  apiError,
			wantErr: prbot.APIError{
				RequestID:  "requestID",
				ErrorCode:  "code",
				Message:    "msg",
				StatusCode: http.StatusBadRequest,
				ErrorText:  errRandom.Error(),
			},
		},
		{
			name: "Should render non APIError",
			err:  errRandom,
			wantErr: prbot.APIError{
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
			prbot.RenderError(rr, req, tt.err)
			assert.Equal(t, tt.wantErr.StatusCode, rr.Code)
			var res prbot.APIError
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
