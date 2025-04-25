package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/stretchr/testify/assert"
)

func TestPullRequestReviewers_GetInputMsg(t *testing.T) {
	ctx := context.Background()
	//nolint:err113 // we don't need static error here for this test
	errRandom := errors.New("random error")

	type args struct {
		ghe             input.GHE
		setExpectations func(d *gh.MockAPI)
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantLogins []string
	}{
		{
			name: "Single user approves once",
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
			wantErr:    false,
			wantLogins: []string{"user1"},
		},
		{
			name: "User submits changes then approves — only latest kept",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return([]*github.PullRequestReview{
							{
								User:  &github.User{Login: github.String("user1")},
								State: github.String("CHANGES_REQUESTED"),
							},
							{
								User:  &github.User{Login: github.String("user1")},
								State: github.String("APPROVED"),
							},
						}, nil)
				},
			},
			wantErr:    false,
			wantLogins: []string{"user1"},
		},
		{
			name: "Error from ListReviews bubbles up",
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

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			// no require, so just assert
			assert.NoError(t, err)
			assert.NotNil(t, got)

			// unmarshal into slice of reviews
			var reviews []*github.PullRequestReview
			err = json.Unmarshal(got, &reviews)
			assert.NoError(t, err, "should unmarshal JSON")

			// collect logins
			var logins []string
			for _, r := range reviews {
				if r.User != nil && r.User.Login != nil {
					logins = append(logins, *r.User.Login)
				}
			}
			// assert we only got each user's most recent APPROVED
			assert.ElementsMatch(t, tt.wantLogins, logins)
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
