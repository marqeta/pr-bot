package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/v50/github"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/stretchr/testify/assert"
)

func mustParse(ts string) time.Time {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		panic(err)
	}
	return t
}

func TestPullRequestReviewers_GetInputMsg(t *testing.T) {
	ctx := context.Background()
	//nolint:goerr113
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
		wantStates map[string]string
	}{
		{
			name: "Single user approves once",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return([]*github.PullRequestReview{
							{
								User:        &github.User{Login: github.String("user1")},
								State:       github.String("APPROVED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T10:00:00Z")},
							},
						}, nil)
				},
			},
			wantErr:    false,
			wantLogins: []string{"user1"},
			wantStates: map[string]string{"user1": "APPROVED"},
		},
		{
			name: "User changes then approves — only latest kept",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return([]*github.PullRequestReview{
							{
								User:        &github.User{Login: github.String("user1")},
								State:       github.String("CHANGES_REQUESTED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T09:00:00Z")},
							},
							{
								User:        &github.User{Login: github.String("user1")},
								State:       github.String("APPROVED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T11:00:00Z")},
							},
						}, nil)
				},
			},
			wantErr:    false,
			wantLogins: []string{"user1"},
			wantStates: map[string]string{"user1": "APPROVED"},
		},
		{
			name: "Multiple users mixed reviews",
			args: args{
				ghe: sampleGHE(),
				setExpectations: func(d *gh.MockAPI) {
					d.EXPECT().ListReviews(ctx, sampleGHE().ToID()).
						Return([]*github.PullRequestReview{
							{
								User:        &github.User{Login: github.String("user1")},
								State:       github.String("APPROVED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T11:00:00Z")},
							},
							{
								User:        &github.User{Login: github.String("user2")},
								State:       github.String("CHANGES_REQUESTED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T12:00:00Z")},
							},
							{
								User:        &github.User{Login: github.String("user1")},
								State:       github.String("CHANGES_REQUESTED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T13:00:00Z")},
							},
							{
								User:        &github.User{Login: github.String("user2")},
								State:       github.String("APPROVED"),
								SubmittedAt: &github.Timestamp{Time: mustParse("2025-04-25T14:00:00Z")},
							},
						}, nil)
				},
			},
			wantErr:    false,
			wantLogins: []string{"user1", "user2"},
			wantStates: map[string]string{"user1": "CHANGES_REQUESTED", "user2": "APPROVED"},
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
			raw, err := prr.GetInputMsg(ctx, tt.args.ghe)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, raw)

			var reviews []*github.PullRequestReview
			assert.NoError(t, json.Unmarshal(raw, &reviews))

			// verify logins and states
			gotLogins := []string{}
			gotStates := map[string]string{}
			for _, r := range reviews {
				if r.User != nil && r.User.Login != nil && r.State != nil {
					login := *r.User.Login
					gotLogins = append(gotLogins, login)
					gotStates[login] = *r.State
				}
			}
			assert.ElementsMatch(t, tt.wantLogins, gotLogins)
			if tt.wantStates != nil {
				assert.Equal(t, tt.wantStates, gotStates)
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
