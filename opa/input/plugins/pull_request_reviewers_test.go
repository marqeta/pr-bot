package plugins

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/stretchr/testify/assert"
)

func TestPullRequestReviewers_GetInputMsg(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
	errRandom := errors.New("random error")

	type args struct {
		ghe             input.GHE
		setExpectations func(d *gh.MockAPI)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should return reviews successfully",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return([]*github.PullRequestReview{
							{
								User:  &github.User{Login: github.String("user1")},
								State: github.String("APPROVED"),
							},
						}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error when ListReviews fails",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return(nil, errRandom)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gh.NewMockAPI(t)
			tt.args.setExpectations(api)
			prr := NewPullRequestReviewers(api)
			got, err := prr.GetInputMsg(ctx, tt.args.ghe)
			if (err != nil) != tt.wantErr {
				t.Errorf("PullRequestReviewers.GetInputMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func sampleGHE() input.GHE {
	return input.GHE{
		Event:  "pull_request",
		Action: "opened",
		PullRequest: &github.PullRequest{
			Number: github.Int(1),
			NodeID: github.String("MDExOlB1bGxSZXF1ZXN0MQ=="),
			User:   &github.User{Login: github.String("octocat")},
		},
		Repository: &github.Repository{
			Owner:    &github.User{Login: github.String("octocat")},
			Name:     github.String("Hello-World"),
			FullName: github.String("octocat/Hello-World"),
		},
	}
}
