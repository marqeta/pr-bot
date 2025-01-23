package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/jonboulle/clockwork"
	"github.com/open-policy-agent/opa/sdk"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/configstore"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/healthcheck"
	"github.com/marqeta/pr-bot/oci"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/input/plugins"
	"github.com/marqeta/pr-bot/pullrequest"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/marqeta/pr-bot/rate"
	"github.com/marqeta/pr-bot/secrets"
	"github.com/marqeta/pr-bot/ui"
	"github.com/marqeta/pr-bot/webhook"
	lim "github.com/mennanov/limiters"
)

const (
	APIURL     = "https://%s/api/v3"
	UploadURL  = "https://%s/api/uploads"
	GraphqlURL = "https://%s/api/graphql"
)

func main() {
	cfg, err := prbot.ParseConfigFiles()
	if err != nil {
		os.Exit(1)
	}

	log.Info().Interface("Config", cfg).Msg("Successfully read config")

	svc, err := prbot.NewService(cfg)
	if err != nil {
		log.Err(err).Msg("Error creating pr-bot service")
		os.Exit(1)
	}

	downloadOCIArtifacts(svc, cfg)

	endpoints := make([]prbot.Endpoint, 0)
	endpoints = append(endpoints, webhookEndpoint(svc, cfg))
	endpoints = append(endpoints, ui.NewEndpoint(svc.EvaluationManager, svc.Metrics))
	svc.MountRoutes(healthcheck.NewEndpoint(svc.Metrics), endpoints)

	srv := prbot.NewServer(cfg, svc.Router)
	go srv.Start()
	srv.WaitForGracefulShutdown()
	svc.Close()
}

func setUpOPAClient(cfg *prbot.Config) client.Client {
	log.Info().Msg("Setting up OPA client")
	config := fmt.Sprintf(`
	{
	   "labels": {
	      "app": "%s",
	      "region": "us-east-2",
	      "environment": "%s"
	   },
	   "bundles": {
	      "local": {
	         "resource": "file:///%s/%s"
	      }
	   }
	}`, cfg.ServiceName, cfg.Env, cfg.OPA.Bundles.Root, cfg.OPA.Bundles.Filename)
	// TODO this can block indefinitely, use channel to signal completion and set timeout
	opaSDK, err := sdk.New(context.Background(), sdk.Options{
		ID:     fmt.Sprintf("%s-%s", cfg.ServiceName, cfg.Env),
		Config: bytes.NewReader([]byte(config)),
	})
	if err != nil {
		log.Err(err).Msg("Error creating OPA SDK client")
		os.Exit(1)
	}
	log.Info().Msg("Successfully created OPA SDK client")
	return client.NewClient(opaSDK)
}

func setUpInputFactory(api gh.API) input.Factory {
	log.Info().Msg("Setting up input factory")
	branchProtection := plugins.NewBranchProtection(api)
	// 100KB size limit
	filesChanged := plugins.NewFilesChanged(api, 100*1000)
	return input.NewFactory(branchProtection, filesChanged)
}

func setUpOPAPolicies(opaClient client.Client) opa.Policy {
	log.Info().Msg("Setting up OPA policies")
	v1 := opa.NewV1Policy(opaClient)
	return opa.NewVersionedPolicy(
		map[string]opa.Policy{"v1": v1},
		opaClient,
	)
}

func FindOPAModules(cfg *prbot.Config) []string {
	log.Info().Msg("Finding OPA modules")
	reader := oci.NewReader()
	filepath := fmt.Sprintf("%s/%s", cfg.OPA.Bundles.Root, cfg.OPA.Bundles.Filename)
	dirs, err := reader.ListDirs(context.Background(), filepath)
	if err != nil {
		log.Err(err).Msg("Error reading OPA bundle directories")
		os.Exit(1)
	}
	modules := reader.FilterModules(context.Background(), dirs)
	log.Info().Interface("Modules", modules).Msg("Found OPA modules")
	return modules
}

func setUpOPAEvaluator(api gh.API, cfg *prbot.Config, svc *prbot.Service) opa.Evaluator {
	log.Info().Msg("Setting up OPA evaluator")
	client := setUpOPAClient(cfg)
	modules := FindOPAModules(cfg)
	policy := setUpOPAPolicies(client)
	factory := setUpInputFactory(api)
	return opa.NewEvaluator(modules, policy, factory, svc.EvaluationManager)
}

func downloadOCIArtifacts(svc *prbot.Service, cfg *prbot.Config) {
	log.Info().Msg("Downloading OPA bundle from ECR")
	puller := oci.NewECRPuller(oci.NewECRCredRetriever(svc.ECR))
	err := puller.Pull(context.Background(), oci.ArtifactID{
		Registry: cfg.OPA.Bundles.ECR.Registry,
		Repo:     cfg.OPA.Bundles.ECR.Repo,
		Tag:      cfg.OPA.Bundles.ECR.Tag,
	}, cfg.OPA.Bundles.Root)
	if err != nil {
		log.Err(err).Msg("Error pulling OPA bundle from ECR")
		os.Exit(1)
	}
	log.Info().Msg("Successfully pulled OPA bundle from ECR")
}

func setupReviewer(svc *prbot.Service, cfg *prbot.Config, api gh.API) review.Reviewer {
	log.Info().Msg("Setting up reviewer")
	// mutex -> dedup -> precond -> rate limited -> reviewer
	base := review.NewReviewer(api, svc.Metrics)
	throttler := setupThrottlers(svc, cfg)
	rateLimited := review.NewRateLimitedReviewer(base, api, throttler)
	precond := review.NewPreCondValidationReviewer(rateLimited)
	dedup := review.NewDedupReviewer(precond, api, cfg.GHE.ServiceAccount)
	locker, err := review.NewLocker(svc.DDB, cfg.Reviewer.Locker.TableName)
	if err != nil {
		log.Err(err).Msg("Error creating review.locker")
		os.Exit(1)
	}
	return review.NewMutexReviewer(dedup, locker)
}

func webhookEndpoint(svc *prbot.Service, cfg *prbot.Config) prbot.Endpoint {
	log.Info().Msg("Setting up webhook endpoint")
	ws, err := svc.Secrets.GetSecret(context.Background(), cfg.AWS.Secrets.Webhook)
	if err != nil {
		log.Err(err).Msg("Error retrieving webhook secret")
		os.Exit(1)
	}
	v3, v4 := setupGHEClients(svc.Secrets, cfg)
	prDao := gh.NewAPI(cfg.Server.Host, cfg.Server.Port, v3, v4, svc.Metrics)
	filter := setupEventFilters(svc, cfg, prDao)
	opaEvaluator := setUpOPAEvaluator(prDao, cfg, svc)
	reviewer := setupReviewer(svc, cfg, prDao)

	handlerV2 := pullrequest.NewEventHandler(opaEvaluator, reviewer, svc.Metrics)
	d := pullrequest.NewDispatcher(handlerV2, filter, svc.Metrics)
	p := webhook.NewGHEventsParser()
	return webhook.NewEndpoint(ws, p, d, svc.Metrics)
}

func setupThrottlers(svc *prbot.Service, cfg *prbot.Config) rate.Throttler {
	author := setupThrottler("AuthorBasedThrottler", rate.AuthorKey, svc, cfg)
	org := setupThrottler("OrgBasedThrottler", rate.OrgKey, svc, cfg)
	repo := setupThrottler("RepoBasedThrottler", rate.RepoKey, svc, cfg)

	return rate.NewFacade(svc.Metrics, repo, org, author)
}

func setupThrottler(name string, keyer rate.Keyer, svc *prbot.Service,
	cfg *prbot.Config) rate.Throttler {

	log.Info().Msgf("Setting up %v Throttler", name)
	cfgStoreName := name + "Config"
	csDao := configstore.NewDynamoDao[*rate.LimiterConfig](svc.DDB, svc.Metrics)
	ticker := clockwork.NewRealClock().NewTicker(cfg.ConfigStore.Refresh)
	configstore, err := configstore.NewDBStore(csDao, cfgStoreName,
		cfg.ConfigStore.Table, ticker, svc.Metrics)
	if err != nil {
		log.Err(err).Msgf("Error reading %v from ddb", cfgStoreName)
		os.Exit(1)
	}

	// all throtllers use the same DDB table
	// but they have different partition key based on rate.keyer
	props := lim.DynamoDBTableProperties{
		TableName:        cfg.Throttler.Table,
		PartitionKeyName: cfg.Throttler.PartitionKey,
		SortKeyName:      cfg.Throttler.SortKey,
		SortKeyUsed:      true,
		TTLFieldName:     cfg.Throttler.TTLFieldName,
	}
	registry := rate.NewSWRegistry(name, svc.DDB, props, clockwork.NewRealClock())
	return rate.NewSlidingWindowLimiter(keyer, registry, configstore)
}

func setupEventFilters(svc *prbot.Service, cfg *prbot.Config, api gh.API) pullrequest.EventFilter {
	log.Info().Msg("Setting up event filters")

	csDao := configstore.NewDynamoDao[*pullrequest.RepoFilterCfg](svc.DDB, svc.Metrics)

	ticker := clockwork.NewRealClock().NewTicker(cfg.ConfigStore.Refresh)
	configstore, err := configstore.NewDBStore(csDao, "RepoFilterConfig",
		cfg.ConfigStore.Table, ticker, svc.Metrics)
	if err != nil {
		log.Err(err).Msg("Error reading RepoFilterConfig from ddb")
		os.Exit(1)
	}
	return pullrequest.NewRepoFilter(configstore, api)
}

func setupGHEClients(sm secrets.Manager, cfg *prbot.Config) (*github.Client, *githubv4.Client) {
	log.Info().Msg("Setting up GHE clients")
	tok, err := sm.GetSecret(context.Background(), cfg.AWS.Secrets.Token)
	if err != nil {
		log.Err(err).Msg("Error retrieving github token")
		os.Exit(1)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok},
	)
	httpClient := oauth2.NewClient(context.Background(), ts)

	a := fmt.Sprintf(APIURL, cfg.GHE.Hostname)
	u := fmt.Sprintf(UploadURL, cfg.GHE.Hostname)
	g := fmt.Sprintf(GraphqlURL, cfg.GHE.Hostname)

	v3, err := github.NewEnterpriseClient(a, u, httpClient)
	if err != nil {
		log.Err(err).Msg("Error creating github client")
		os.Exit(1)
	}
	v4 := githubv4.NewEnterpriseClient(g, httpClient)
	return v3, v4
}
