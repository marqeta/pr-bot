package data

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	pe "github.com/marqeta/pr-bot/errors"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	Message    string `json:"message,omitempty"`
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (c *controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oplog := httplog.LogEntry(ctx)
	reqID := middleware.GetReqID(ctx)

	// do auth verification
	// TODO SigV4 verification for presigned get-caller-identity request
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
	} else {
		render.Render(w, r, &Response{
			StatusCode: http.StatusOK,
			RequestID:  reqID,
			Message:    "payload stored successfully",
		})
	}
}
