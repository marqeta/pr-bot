package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type APIError struct {
	RequestID  string `json:"request_id"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code"`
	// JSON doesnt marshal Error type; store the error text as a string field
	ErrorText string `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%d - %s : %s : %s", e.StatusCode, e.ErrorCode, e.Message, e.ErrorText)
}

func (e APIError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (e APIError) IsUserError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

func (e APIError) IsServiceError() bool {
	return !e.IsUserError()
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	var ae APIError
	if errors.As(err, &ae) {
		_ = render.Render(w, r, ae)
		return
	}
	_ = render.Render(w, r, APIError{
		RequestID:  middleware.GetReqID(r.Context()),
		ErrorCode:  "Unknown",
		StatusCode: http.StatusInternalServerError,
		ErrorText:  err.Error(),
	})
}

func TooManyRequestError(ctx context.Context, msg string, err error) APIError {
	return APIError{
		RequestID:  middleware.GetReqID(ctx),
		ErrorCode:  "TooManyRequestsError",
		Message:    msg,
		StatusCode: http.StatusTooManyRequests,
		ErrorText:  err.Error(),
	}
}

func InValidRequestError(ctx context.Context, msg string, err error) APIError {
	return APIError{
		RequestID:  middleware.GetReqID(ctx),
		ErrorCode:  "InValidRequestError",
		Message:    msg,
		StatusCode: http.StatusBadRequest,
		ErrorText:  err.Error(),
	}
}

func ServiceFault(ctx context.Context, msg string, err error) APIError {
	return APIError{
		RequestID:  middleware.GetReqID(ctx),
		ErrorCode:  "ServiceFault",
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
		ErrorText:  err.Error(),
	}
}

func UserError(ctx context.Context, msg string, err error) APIError {
	return APIError{
		RequestID:  middleware.GetReqID(ctx),
		ErrorCode:  "UserError",
		Message:    msg,
		StatusCode: http.StatusBadRequest,
		ErrorText:  err.Error(),
	}
}
