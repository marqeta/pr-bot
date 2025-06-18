package data

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"

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

type GetCallerIdentityResponse struct {
	XMLName xml.Name `xml:"GetCallerIdentityResponse"`
	Result  struct {
		Arn     string `xml:"Arn"`
		Account string `xml:"Account"`
		UserID  string `xml:"UserId"`
	} `xml:"GetCallerIdentityResult"`
}

var HTTPGet = http.Get

func (e *Response) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (c *controller) HandleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	ctx = evaluation.SetDeliveryID(r.Context(), uuid.NewString())
	oplog := httplog.LogEntry(ctx)

	callerArn, err := verifySTSIdentity(ctx, r)
	if err != nil {
		oplog.Err(err).Msg("identity verification failed")
		pe.RenderError(w, r, err)
		return
	}

	oplog.Info().Str("callerArn", callerArn).Msg("identity verified")

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

func verifySTSIdentity(ctx context.Context, r *http.Request) (string, error) {
	stsURL := r.Header.Get("X-Aws-Sts-Signature")
	if stsURL == "" {
		return "", pe.UserError(ctx, "missing STS signature header", nil)
	}

	resp, err := HTTPGet(stsURL)
	if err != nil {
		return "", pe.UserError(ctx, "failed to call STS", err)
	}
	defer resp.Body.Close()

	oplog := httplog.LogEntry(ctx)

	oplog.Info().Int("status", resp.StatusCode).Str("url", stsURL).Msg("STS GetCallerIdentity response")

	if resp.StatusCode != http.StatusOK {
		//nolint:goerr113
		return "", pe.UserError(ctx, "STS responded with non-200", errors.New("non-200 response from STS"))
	}

	// Log the actual status code

	stsBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", pe.UserError(ctx, "STS read error", err)
	}

	var stsResp GetCallerIdentityResponse
	if err := xml.Unmarshal(stsBody, &stsResp); err != nil {
		return "", pe.UserError(ctx, "STS parse error", err)
	}

	callerArn := stsResp.Result.Arn
	oplog.Info().Str("callerarn", callerArn).Msg("STS CALLER ARN")

	if callerArn == "" || !strings.Contains(callerArn, ":assumed-role/s--polynator") {
		//nolint:goerr113
		return "", pe.UserError(ctx, "unauthorized identity", errors.New("unauthorized identity"))
	}

	return callerArn, nil
}
