package data

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/opa/evaluation"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	Message    string `json:"message,omitempty"`
}

func (e *Response) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (c *controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	ctx = evaluation.SetDeliveryID(r.Context(), uuid.NewString())
	oplog := httplog.LogEntry(ctx)

	callerArn, err := c.verifier.Verify(ctx, r)
	if err != nil {
		oplog.Err(err).Msg("identity verification failed")
		pe.RenderError(w, r, err)
		return
	}

	oplog.Info().Str("sigv4 callerArn", callerArn).Msg("identity verified")

	metadata, err := c.dao.ToMetadata(ctx, r)
	if err != nil {
		oplog.Err(err).Msg("error parsing metadata")
		pe.RenderError(w, r, pe.UserError(ctx, "error parsing metadata", err))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		oplog.Err(err).Msg("error reading request body")
		pe.RenderError(w, r, pe.UserError(ctx, "error reading request body", err))
		return
	}
	defer r.Body.Close()

	err = c.dao.StorePayload(ctx, metadata, body)
	if err != nil {
		oplog.Err(err).Msg("error storing payload")
		pe.RenderError(w, r, err)
		return
	}

	err = c.handler.EvalAndReviewDataEvent(ctx, metadata)
	if err != nil {
		oplog.Err(err).Msg("error evaluating policies during data event")
		pe.RenderError(w, r, err)
		return
	}

	_ = render.Render(w, r, &Response{
		StatusCode: http.StatusOK,
		RequestID:  reqID,
		Message:    "payload stored successfully",
	})
}
