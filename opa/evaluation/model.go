package evaluation

import (
	"context"

	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
)

// Key to use when setting the delivery ID.
type ctxKeyDeliveryID int

const DeliveryIDKey ctxKeyDeliveryID = 0

type Result struct {
	Result types.Result `json:"result"`
	Err    error        `json:"err"`
}

type Report struct {
	ReportMetadata
	Breakdown map[string]Result `json:"breakdown"`
	Input     *input.Model      `json:"input"`
}

type ReportMetadata struct {
	PR            string `json:"pr"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	RequestID     string `json:"request_id"`
	DeliveryID    string `json:"delivery_id"`
	PolicyVersion string `json:"policy_version"`
	Outcome       Result `json:"outcome"`
	ExpireAt      int64  `json:"expire_at"`
	CreatedAt     int64  `json:"created_at"`
	Event         string `json:"event"`
	Action        string `json:"action"`
}

func GetDeliveryID(ctx context.Context) string {
	if val, ok := ctx.Value(DeliveryIDKey).(string); ok {
		return val
	}
	return ""
}
