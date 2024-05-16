package rate_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/configstore"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/rate"

	lim "github.com/mennanov/limiters"
)

func Test_swLimiter_ShouldThrottle(t *testing.T) {
	//nolint:goerr113
	randErr := errors.New("random error")
	ctx := context.Background()
	type fields struct {
		keyer rate.Keyer
		cfg   *rate.LimiterConfig
	}
	type args struct {
		id id.PR
	}
	tests := []struct {
		name            string
		fields          fields
		setExpectations func(reg *rate.MockGetter, limiter *rate.MockLimiter)
		args            args
		wantErr         error
	}{
		{
			name: "Should return 429 when throttled for repo throttle",
			fields: fields{
				keyer: rate.RepoKey,
				cfg: &rate.LimiterConfig{
					Default: limit(5, 10),
				}},
			setExpectations: func(reg *rate.MockGetter, limiter *rate.MockLimiter) {
				reg.EXPECT().Name().Return("Mock").Once()
				reg.EXPECT().GetOrCreate(ctx, "Repo/owner1/repo1", limit(5, 10)).
					Return(limiter, nil).Once()
				limiter.EXPECT().Limit(ctx).
					Return(5*time.Second, lim.ErrLimitExhausted).Once()
			},
			args: args{
				id: ID("owner1", "repo1", "author1"),
			},
			wantErr: prbot.TooManyRequestError(ctx,
				fmt.Sprintf("%v throttled request for key %v, try again in %v",
					"Mock", "Repo/owner1/repo1", 5*time.Second),
				lim.ErrLimitExhausted),
		},
		{
			name: "Should return 429 when throttled for org throttle",
			fields: fields{
				keyer: rate.OrgKey,
				cfg: &rate.LimiterConfig{
					Default: limit(5, 10),
				}},
			setExpectations: func(reg *rate.MockGetter, limiter *rate.MockLimiter) {
				reg.EXPECT().Name().Return("Mock").Once()
				reg.EXPECT().GetOrCreate(ctx, "Org/owner1", limit(5, 10)).
					Return(limiter, nil).Once()
				limiter.EXPECT().Limit(ctx).
					Return(5*time.Second, lim.ErrLimitExhausted).Once()
			},
			args: args{
				id: ID("owner1", "repo1", "author1"),
			},
			wantErr: prbot.TooManyRequestError(ctx,
				fmt.Sprintf("%v throttled request for key %v, try again in %v",
					"Mock", "Org/owner1", 5*time.Second),
				lim.ErrLimitExhausted),
		},
		{
			name: "Should return 429 when throttled for author throttle",
			fields: fields{
				keyer: rate.AuthorKey,
				cfg: &rate.LimiterConfig{
					Default: limit(5, 10),
				}},
			setExpectations: func(reg *rate.MockGetter, limiter *rate.MockLimiter) {
				reg.EXPECT().Name().Return("Mock").Once()
				reg.EXPECT().GetOrCreate(ctx, "Author/author1", limit(5, 10)).
					Return(limiter, nil).Once()
				limiter.EXPECT().Limit(ctx).
					Return(5*time.Second, lim.ErrLimitExhausted).Once()
			},
			args: args{
				id: ID("owner1", "repo1", "author1"),
			},
			wantErr: prbot.TooManyRequestError(ctx,
				fmt.Sprintf("%v throttled request for key %v, try again in %v",
					"Mock", "Author/author1", 5*time.Second),
				lim.ErrLimitExhausted),
		},
		{
			name: "Should error from registry",
			fields: fields{
				keyer: rate.RepoKey,
				cfg: &rate.LimiterConfig{
					Default: limit(6, 10),
				}},
			setExpectations: func(reg *rate.MockGetter, _ *rate.MockLimiter) {
				reg.EXPECT().GetOrCreate(ctx, "Repo/owner1/repo1", limit(6, 10)).
					Return(nil, randErr).Once()
			},
			args: args{
				id: ID("owner1", "repo1", "author1"),
			},
			wantErr: randErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockReg := rate.NewMockGetter(t)

			store, err := configstore.NewInMemoryStore(tt.fields.cfg)
			assert.Nil(t, err)

			limiter := rate.NewMockLimiter(t)

			tt.setExpectations(mockReg, limiter)

			sw := rate.NewSlidingWindowLimiter(tt.fields.keyer, mockReg, store)

			err = sw.ShouldThrottle(ctx, tt.args.id)
			if !assert.EqualError(t, err, tt.wantErr.Error()) {
				t.Errorf("swLimiter.ShouldThrottle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ID(owner, repo, author string) id.PR {

	return id.PR{
		Owner:        owner,
		Repo:         repo,
		RepoFullName: fmt.Sprintf("%s/%s", owner, repo),
		Author:       author,
	}
}
