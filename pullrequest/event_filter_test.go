package pullrequest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/marqeta/pr-bot/configstore"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/pullrequest"
)

func TestRepoFilterCfg_Update(t *testing.T) {
	type args struct {
		a []string
		d []string
		i []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should create event filter with allow list",
			args: args{
				a: []string{"a/.+"},
			},
			wantErr: false,
		},
		{
			name: "Should create event filter with deny list",
			args: args{
				d: []string{"a/.+"},
			},
			wantErr: false,
		},
		{
			name: "Should create event filter with allow and deny list",
			args: args{
				d: []string{"a/.+", "b/.+"},
				a: []string{"b/.+", "c/.+", "def"},
			},
			wantErr: false,
		},
		{
			name: "Should create event filter with allow, deny list and ignore topics",
			args: args{
				d: []string{"a/.+", "b/.+"},
				a: []string{"b/.+", "c/.+", "def"},
				i: []string{"topic", "topic2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &pullrequest.RepoFilterCfg{
				Allowlist:    tt.args.a,
				Denylist:     tt.args.d,
				IgnoreTopics: tt.args.i,
			}
			err := cfg.Update()
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoFilterCfg.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, cfg.AllowlistRegex, len(cfg.Allowlist))
			assert.Len(t, cfg.DenylistRegex, len(cfg.DenylistRegex))
			assert.Len(t, cfg.IgnoreTopicsMap, len(cfg.IgnoreTopics))
		})
	}
}

func Test_repoNameFilter_ShouldHandle(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		cfg *pullrequest.RepoFilterCfg
	}
	type args struct {
		id id.PR
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		setExpectations func(dao *gh.MockAPI, id id.PR)
		want            bool
		wantErr         bool
	}{
		{
			name: "ShouldHandle=false when not matching both allow list or deny list are empty",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{},
					Denylist:  []string{},
				},
			},
			args: args{
				id: idWithFullName("a/c"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, nil).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false when not matching both allow list or deny list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"def/312"},
				},
			},
			args: args{
				id: idWithFullName("a/c"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, nil).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false when matching deny list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"def/312"},
				},
			},
			args: args{
				id: idWithFullName("def/312"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=true when matching allow list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"def/312"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false when deny list takes precedence",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"abc/123"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false with atleast one matching deny list entry",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"def/312", "zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("def/312"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false with partial match in deny list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/123"},
					Denylist:  []string{"def/*", "zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("def/312"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=true with partial match in allow list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/*"},
					Denylist:  []string{"def/*", "zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ShouldHandle=true when allowlist=.+",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{".+"},
					Denylist:  []string{"def/*", "zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false when allowlist=.+ but matching denylist",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{".+"},
					Denylist:  []string{"def/*", "zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("def/123"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false if dao.ListAllTopics error",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{".+"},
					Denylist:  []string{"zxc/644"},
				},
			},
			args: args{
				id: idWithFullName("def/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{}, errRandom).Once()
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "ShouldHandle=false if repo has pr-bot-ignore",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist:    []string{".+"},
					Denylist:     []string{"zxc/644"},
					IgnoreTopics: []string{"pr-bot-ignore"},
				},
			},
			args: args{
				id: idWithFullName("def/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).
					Return([]string{"random", "pr-bot-ignore", "abc"}, nil).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false if repo has ignore topic and in allow list",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist:    []string{"def/123"},
					Denylist:     []string{"zxc/644"},
					IgnoreTopics: []string{"pr-bot-ignore"},
				},
			},
			args: args{
				id: idWithFullName("def/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return([]string{"random", "pr-bot-ignore", "random"}, nil).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=true for staging repo in staging env",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc-staging/.+"},
					Denylist:  []string{"abc/.+"},
				},
			},
			args: args{
				id: idWithFullName("abc-staging/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return(nil, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false for prod repo in staging env",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc-staging/.+"},
					Denylist:  []string{"abc/.+"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=false for staging repo in prod env",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/.+"},
					Denylist:  []string{"abc-staging/.+"},
				},
			},
			args: args{
				id: idWithFullName("abc-staging/123"),
			},
			setExpectations: func(_ *gh.MockAPI, _ id.PR) {
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ShouldHandle=true for prod repo in prod env",
			fields: fields{
				cfg: &pullrequest.RepoFilterCfg{
					Allowlist: []string{"abc/.+"},
					Denylist:  []string{"abc-staging/.+"},
				},
			},
			args: args{
				id: idWithFullName("abc/123"),
			},
			setExpectations: func(dao *gh.MockAPI, id id.PR) {
				dao.EXPECT().ListAllTopics(ctx, id).Return(nil, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := gh.NewMockAPI(t)
			store, err := configstore.NewInMemoryStore(tt.fields.cfg)
			if err != nil {
				t.Errorf("error creating config store with config %v", tt.fields.cfg)
			}
			f := pullrequest.NewRepoFilter(store, dao)
			tt.setExpectations(dao, tt.args.id)
			got, err := f.ShouldHandle(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("repoNameFilter.ShouldHandle().error = %v, want %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repoNameFilter.ShouldHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func idWithFullName(fullname string) id.PR {

	return id.PR{
		Owner:        "owner1",
		Repo:         "repo1",
		Number:       1,
		NodeID:       "nodeid1",
		RepoFullName: fullname,
	}
}
