package prbot

import (
	"context"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/secrets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Service struct {
	Router            *chi.Mux
	Config            *Config
	Logger            zerolog.Logger
	Secrets           secrets.Manager
	DDB               *dynamodb.Client
	ECR               *ecr.Client
	Metrics           metrics.Emitter
	EvaluationManager evaluation.Manager
}

func NewService(cfg *Config) (*Service, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.AWS.Region))
	if err != nil {
		log.Err(err).Msg("unable to load SDK config")
		return nil, err
	}
	metrics, err := setupMetricsEmitter(cfg)
	if err != nil {
		log.Err(err).Msg("unable to create datadog statsd metrics emitter")
		return nil, err
	}

	ddb := dynamodb.NewFromConfig(awsConfig)
	manager := evaluation.NewManager(ddb, cfg.OPA.Bundles.ECR.Tag,
		cfg.OPA.EvaluationReport.TTL, metrics, cfg.OPA.EvaluationReport.TableName)
	s := &Service{
		Router:            chi.NewRouter(),
		Config:            cfg,
		Logger:            setupLogger(cfg),
		Secrets:           secrets.NewManager(secretsmanager.NewFromConfig(awsConfig)),
		DDB:               ddb,
		ECR:               ecr.NewFromConfig(awsConfig),
		Metrics:           metrics,
		EvaluationManager: manager,
	}
	logAWSIdentity(awsConfig)
	return s, nil
}

func logAWSIdentity(awsConfig aws.Config) {
	stsClient := sts.NewFromConfig(awsConfig)
	caller, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Err(err).Msg("GetCallerIdentity failed")
		return
	}
	log.Info().Interface("Caller", caller).Msg("AWS identity")
}

func setupLogger(cfg *Config) zerolog.Logger {
	isLocal := false
	if cfg.Env == "local" {
		isLocal = true
	}
	logger := httplog.NewLogger(cfg.ServiceName, httplog.Options{
		Concise: isLocal,
		JSON:    !isLocal,
		// service, env, version tag is used for unifed service logging in datadog
		Tags: map[string]string{
			"env":     cfg.Env,
			"version": cfg.Build.Version,
			// service tag is added automatically
		},
	})
	// Adds caller to the logs and make httplog the default logger
	log.Logger = logger.With().Caller().Logger()
	return log.Logger
}

func setupMetricsEmitter(cfg *Config) (metrics.Emitter, error) {
	if cfg.Env == "local" {
		return metrics.NewNoopEmitter(), nil
	}
	podName, err := os.Hostname()
	if err != nil {
		log.Err(err).Msg("Error getting pod name")
		return nil, err
	}

	// uses DD_AGENT_HOST env var to find the statsd ip and uses default port of 8125
	// service, env, version tags are already added by default from env vars
	// {env_var, tag_name}
	// {"DD_ENV", "env"},                         // The name of the env in which the service runs.
	// {"DD_SERVICE", "service"},                 // The name of the running service.
	// {"DD_VERSION", "version"},                 // The current version of the running service.
	client, err := statsd.New("", statsd.WithNamespace("pr-bot"),
		statsd.WithoutTelemetry(), statsd.WithoutDevMode(),
		statsd.WithTags([]string{"pod_name:" + podName}))
	if err != nil {
		log.Err(err).Msg("Error creating datadog client")
		return nil, err
	}
	return metrics.NewEmitter(client), nil
}

func (s *Service) MountRoutes(healthcheck Endpoint, endpoints []Endpoint) {
	s.Router.Use(middleware.Timeout(60 * time.Second))
	// Request logger adds request id handler and recovery handler
	// skip logging health check rqeuests
	s.Router.Use(httplog.RequestLogger(s.Logger, []string{healthcheck.Path()}))

	s.Router.Mount(healthcheck.Path(), healthcheck.Routes())

	for _, endpoint := range endpoints {
		s.Router.Mount(endpoint.Path(), endpoint.Routes())
	}
}

func (s *Service) Close() {
	s.Metrics.Close()
}
