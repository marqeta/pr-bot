package evaluation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
)

func Test_reportBuilder(t *testing.T) {
	ctx := context.TODO()
	dao := evaluation.NewMockDao(t)
	noop := metrics.NewNoopEmitter()
	//nolint:goerr113
	randomErr := fmt.Errorf("random error")
	now := time.Now()
	m := evaluation.NewManager(dao, "1.0.-163", 5*time.Hour, noop, "tableName")
	b := m.NewReportBuilder(ctx, "ci/pr-bot/123", "1234567890", "deliveryID")
	assert.NotNil(t, b)

	b.SetInput(randomModel())
	b.AddModuleResult("module1", approve())
	b.AddModuleResult("module2", reqChanges())
	b.AddModuleResult("module3", err(randomErr))
	b.SetOutcome(err(randomErr))
	report := b.GetReport()

	assert.Equal(t, "ci/pr-bot/123", report.PR)
	assert.Equal(t, "1234567890", report.RequestID)
	assert.Equal(t, "deliveryID", report.DeliveryID)
	assert.Equal(t, "Automatic Dockerfile Image Updater", report.Title)
	assert.Equal(t, "svc-ci-dfiu", report.Author)
	assert.Equal(t, "1.0.-163", report.PolicyVersion)
	assert.Equal(t, now.Add(5*time.Hour).Unix(), report.ExpireAt)
	assert.Equal(t, randomModel().Event, report.Event)
	assert.Equal(t, randomModel().Action, report.Action)
	assert.Equal(t, randomModel(), report.Input)
	assert.Equal(t, 3, len(report.Breakdown))
	assert.Equal(t, approve(), report.Breakdown["module1"])
	assert.Equal(t, reqChanges(), report.Breakdown["module2"])
	assert.Equal(t, err(randomErr), report.Breakdown["module3"])
	assert.Equal(t, err(randomErr), report.Outcome)
}

func randomModel() *input.Model {
	return &input.Model{
		Event:  "pull_request",
		Action: "opened",
		PullRequest: &github.PullRequest{
			Number: aws.Int(259),
			Title:  aws.String("Automatic Dockerfile Image Updater"),
			User: &github.User{
				Login: aws.String("svc-ci-dfiu"),
			},
		},
		Repository: &github.Repository{
			Name:     aws.String("terraform-provider-oci"),
			Owner:    &github.User{Login: aws.String("ci")},
			FullName: aws.String("ci/terraform-provider-oci"),
		},
	}
}

func approve() evaluation.Result {
	return evaluation.Result{
		Result: types.Result{
			Track: true,
			Review: types.Review{
				Type: types.Approve,
				Body: "LGTM",
			},
		},
		Err: nil,
	}

}

func reqChanges() evaluation.Result {
	return evaluation.Result{
		Result: types.Result{
			Track: true,
			Review: types.Review{
				Type: types.RequestChanges,
				Body: "need unit tests",
			},
		},
		Err: nil,
	}
}

func err(e error) evaluation.Result {
	return evaluation.Result{
		Err: e,
	}
}
