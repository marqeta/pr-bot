package webhook

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/google/go-github/v50/github"

	pe "github.com/marqeta/pr-bot/errors"
	"github.com/marqeta/pr-bot/opa/evaluation"
)

type EventResponse struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	DeliveryID string `json:"delivery_id,omitempty"`
}

func (e *EventResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (c *controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oplog := httplog.LogEntry(ctx)
	reqID := middleware.GetReqID(ctx)

	// Validate payload using shared webhook secret
	payload, err := c.parser.ValidatePayload(r, []byte(c.webhookSecret))
	if err != nil {
		oplog.Err(err).Msg("could not validate webhook payload")
		pe.RenderError(w, r,
			pe.InValidRequestError(ctx, "could not validate webhook payload", err))
		return
	}

	// Parse event payload
	event, err := c.parser.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		oplog.Err(err).Msg("could not parse webhook")
		pe.RenderError(w, r,
			pe.InValidRequestError(ctx, "could not parse webhook", err))
		return
	}

	deliveryID := github.DeliveryID(r)
	httplog.LogEntrySetField(ctx, "deliveryID", deliveryID)
	ctx = context.WithValue(ctx, evaluation.DeliveryIDKey, deliveryID)

	eventName := github.WebHookType(r)
	httplog.LogEntrySetField(ctx, "eventName", eventName)
	oplog = httplog.LogEntry(ctx)

	switch event := event.(type) {

	case *github.PullRequestEvent:
		err = c.dispatcher.Dispatch(ctx, deliveryID, eventName, event)
	case *github.PullRequestReviewEvent:
		err = c.dispatcher.DispatchReview(ctx, deliveryID, eventName, event)

	default:
		oplog.Info().Msgf("No Handlers registered for Event: %s", eventName)
	}

	if err != nil {
		oplog.Err(err).Msg("Error Handling Event")
		pe.RenderError(w, r, err)
	} else {
		_ = render.Render(w, r, &EventResponse{
			StatusCode: http.StatusAccepted,
			RequestID:  reqID,
			DeliveryID: deliveryID,
		})
	}
}
