package identity

import (
	"context"

	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/httplog"
	pe "github.com/marqeta/pr-bot/errors"
	"io"
	"net/http"
	"net/url"
)

type IdentityFetcher interface {
	FetchCallerIdentity(ctx context.Context, r *http.Request) (*CallerIdentity, error)
}

type DefaultFetcher struct {
	HTTPClient *http.Client
}


func (f *DefaultFetcher) FetchCallerIdentity(ctx context.Context, r *http.Request) (*CallerIdentity, error) {
	oplog := httplog.LogEntry(ctx)
	rawPresignedURL := r.Header.Get("X-AWS-STS-SIGNATURE")
	if rawPresignedURL == "" {
		return nil, pe.UserError(ctx, "missing STS signature header", errors.New("no signature provided"))
	}

	parsedURL, err := url.Parse(rawPresignedURL)
	if err != nil {
		return nil, pe.UserError(ctx, "invalid STS presigned URL", err)
	}

	reconstructedURL := fmt.Sprintf("https://sts.amazonaws.com/?%s", parsedURL.RawQuery)
	secureSTSURL, err := url.Parse(reconstructedURL)
	if err != nil {
		return nil, pe.UserError(ctx, "could not build secure STS URL", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", secureSTSURL.String(), nil)
	if err != nil {
		return nil, pe.UserError(ctx, "failed to construct STS request", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.HTTPClient.Do(req)
	if err != nil {
		return nil, pe.UserError(ctx, "failed to call STS", err)
	}
	defer resp.Body.Close()

	oplog.Info().Int("status", resp.StatusCode).Str("url", secureSTSURL.String()).Msg("STS GetCallerIdentity response")

	if resp.StatusCode != http.StatusOK {
		return nil, pe.UserError(ctx, "STS responded with non-200", errors.New("non-200 response from STS"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pe.UserError(ctx, "STS read error", err)
	}
	oplog.Info().Str("raw_sts_body", string(body)).Msg("STS raw response")

	var jsonResp GetCallerIdentityJSON
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, pe.UserError(ctx, "STS JSON parse error", err)
	}

	return &jsonResp.GetCallerIdentityResponse.GetCallerIdentityResult, nil
}
