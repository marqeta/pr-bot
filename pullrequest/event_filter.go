package pullrequest

import (
	"context"
	"regexp"

	"github.com/go-chi/httplog"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/configstore"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
)

//go:generate mockery --name EventFilter --testonly
type EventFilter interface {
	ShouldHandle(ctx context.Context, id id.PR) (bool, error)
}

type RepoFilterCfg struct {
	Allowlist       []string         `dynamodbav:"allowlist"`
	Denylist        []string         `dynamodbav:"denylist"`
	IgnoreTopics    []string         `dynamodbav:"ignore_topics"`
	AllowlistRegex  []*regexp.Regexp `dynamodbav:"-"`
	DenylistRegex   []*regexp.Regexp `dynamodbav:"-"`
	IgnoreTopicsMap map[string]bool  `dynamodbav:"-"`
}

// Update hook will be called everytime configstore reads an updated value.
// It is used to populate regex and maps
func (cfg *RepoFilterCfg) Update() error {

	allow := make([]*regexp.Regexp, 0)
	for _, r := range cfg.Allowlist {
		reg, err := regexp.Compile(r)
		if err != nil {
			return err
		}
		allow = append(allow, reg)
	}

	deny := make([]*regexp.Regexp, 0)
	for _, r := range cfg.Denylist {
		reg, err := regexp.Compile(r)
		if err != nil {
			return err
		}
		deny = append(deny, reg)
	}
	cfg.AllowlistRegex = allow
	cfg.DenylistRegex = deny
	cfg.IgnoreTopicsMap = make(map[string]bool)

	for _, topic := range cfg.IgnoreTopics {
		cfg.IgnoreTopicsMap[topic] = true
	}

	return nil
}

type repoFilter struct {
	cfgStore configstore.Getter[*RepoFilterCfg]
	dao      gh.API
}

func NewRepoFilter(store configstore.Getter[*RepoFilterCfg], dao gh.API) EventFilter {
	return &repoFilter{
		cfgStore: store,
		dao:      dao,
	}
}

// Filter implements Filter
func (f *repoFilter) ShouldHandle(ctx context.Context, id id.PR) (bool, error) {

	oplog := httplog.LogEntry(ctx)
	cfg, err := f.cfgStore.Get()
	if err != nil {
		oplog.Err(err).Msgf("error while retrieveing repo filter cfg")
		return false, err
	}
	if !isEmpty(cfg.DenylistRegex) && isMatching(cfg.DenylistRegex, id.RepoFullName) {
		oplog.Info().Msgf("ShouldHandle=false %v is matching the repo deny list", id.RepoFullName)
		return false, nil
	}

	hasIgnoreTopic, err := f.hasIgnoreTopic(ctx, id, cfg)
	if err != nil {
		oplog.Err(err).Msgf("ShouldHandle=false error while listing topics on %v", id.RepoFullName)
		return false, err
	}

	if hasIgnoreTopic {
		oplog.Info().Msgf("ShouldHandle=false %v has one of %v topic set", id.RepoFullName, cfg.IgnoreTopics)
		return false, nil
	}

	if !isEmpty(cfg.AllowlistRegex) && isMatching(cfg.AllowlistRegex, id.RepoFullName) {
		oplog.Info().Msgf("ShouldHandle=true %v is matching the repo allow list", id.RepoFullName)
		return true, nil
	}

	oplog.Info().Msgf("ShouldHandle=false %v is not matching the repo allow or deny list", id.RepoFullName)
	return false, nil
}

func (f *repoFilter) hasIgnoreTopic(ctx context.Context, id id.PR, cfg *RepoFilterCfg) (bool, error) {
	topics, err := f.dao.ListAllTopics(ctx, id)
	if err != nil {
		return true, prbot.ServiceFault(ctx, "error listing topics on repo", err)
	}
	for _, topic := range topics {
		if cfg.IgnoreTopicsMap[topic] {
			return true, nil
		}
	}
	return false, nil
}

func isMatching(regexs []*regexp.Regexp, name string) bool {
	for _, r := range regexs {
		if r.MatchString(name) {
			return true
		}
	}
	return false
}

func isEmpty(s []*regexp.Regexp) bool {
	return len(s) == 0
}
