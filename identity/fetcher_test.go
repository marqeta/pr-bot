package identity_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marqeta/pr-bot/identity"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestFetchCallerIdentity(t *testing.T) {
	validSTSResponse := `{
		"GetCallerIdentityResponse": {
			"GetCallerIdentityResult": {
				"Arn": "arn:aws:sts::123456789012:assumed-role/MyRoleName/MySessionName",
				"Account": "123456789012",
				"UserId": "AIDABCDEFGHIJKL12345"
			}
		}
	}`

	tests := []struct {
		name               string
		sigHeader          string
		mockResponse       *http.Response
		mockError          error
		expectedErr        string
		expectedCallerArn  string
	}{
		{
			name:        "missing signature header",
			sigHeader:   "",
			expectedErr: "missing STS signature header",
		},
		{
			name:        "invalid presigned URL",
			sigHeader:   ":::invalid-url:::",
			expectedErr: "invalid STS presigned URL",
		},
		{
			name:        "error calling STS",
			sigHeader:   "Action=GetCallerIdentity&Version=2011-06-15",
			mockError:   errors.New("http client error"),
			expectedErr: "failed to call STS",
		},
		{
			name:      "non-200 STS response",
			sigHeader: "Action=GetCallerIdentity&Version=2011-06-15",
			mockResponse: &http.Response{
				StatusCode: 403,
				Body:       io.NopCloser(strings.NewReader("Forbidden")),
			},
			expectedErr: "STS responded with non-200",
		},
		{
			name:      "malformed JSON response",
			sigHeader: "Action=GetCallerIdentity&Version=2011-06-15",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("{not-json")),
			},
			expectedErr: "STS JSON parse error",
		},
		{
			name:      "valid STS response",
			sigHeader: "Action=GetCallerIdentity&Version=2011-06-15",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(validSTSResponse))),
			},
			expectedCallerArn: "arn:aws:sts::123456789012:assumed-role/MyRoleName/MySessionName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return tt.mockResponse, tt.mockError
				},
			}
			client := &http.Client{Transport: mockTransport}
			fetcher := &identity.DefaultFetcher{
				HTTPClient: client,
			}

			req := httptest.NewRequest("GET", "/", nil)
			if tt.sigHeader != "" {
				req.Header.Set("X-AWS-STS-SIGNATURE", tt.sigHeader)
			}

			ctx := context.Background()
			caller, err := fetcher.FetchCallerIdentity(ctx, req)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCallerArn, caller.Arn)
			}
		})
	}
}
