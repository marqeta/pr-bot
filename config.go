package prbot

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Build struct {
		Commit  string `yaml:"Commit" env:"COMMIT" env-description:"Commit hash of the build passed"`
		Repo    string `yaml:"Repo" env:"REPO" env-description:"Github repo containing the source of the build"`
		Version string `yaml:"Version" env:"VERSION" env-description:"Version number 1-0.14.qwehsadhascxzva"`
	} `yaml:"Build" env-prefix:"BUILD_"`
	ServiceName string `yaml:"ServiceName" env:"SERVICE_NAME" env-default:"pr-bot" env-description:"Name of the service"`
	Env         string `yaml:"Env" env:"ENV" env-default:"local" env-description:"stage name; one of local|dev|staging|prod"`
	Server      struct {
		Host string `yaml:"Host" env:"HOST" env-description:"Server host"`
		Port int    `yaml:"Port" env:"PORT" env-description:"Server port"`
		TLS  struct {
			CertFile string `yaml:"CertFile" env:"CERT_FILE" env-description:"Server certificate file path"`
			KeyFile  string `yaml:"KeyFile" env:"KEY_FILE" env-description:"File path of server cert's private key"`
		} `yaml:"TLS" env-prefix:"TLS_"`
	} `yaml:"Server" env-prefix:"SERVER_"`
	AWS struct {
		Region  string `yaml:"REGION" env:"REGION" env-default:"us-east-1"`
		Secrets struct {
			Webhook string `yaml:"Webhook" env:"WEBHOOK" env-default:"/ci/pr-bot/webhook"`
			Token   string `yaml:"Token" env:"TOKEN" env-default:"/ci/pr-bot/token"`
		} `yaml:"Secrets" env-prefix:"SECRETS_"`
	} `yaml:"AWS" env-prefix:"AWS_"`
	OPA struct {
		Bundles struct {
			Root     string `yaml:"Root" env:"ROOT"`
			Filename string `yaml:"Filename" env:"FILENAME"`
			ECR      struct {
				Registry string `yaml:"Registry" env:"REGISTRY"`
				Repo     string `yaml:"Repo" env:"REPO"`
				Tag      string `yaml:"Tag" env:"TAG"`
			} `yaml:"ECR" env-prefix:"ECR_"`
		} `yaml:"Bundles" env-prefix:"BUNDLES_"`
		EvaluationReport struct {
			TTL       time.Duration `yaml:"TTL" env:"TTL"`
			TableName string        `yaml:"TableName" env:"TABLE_NAME"`
		} `yaml:"EvaluationReport" env-prefix:"EVALUATION_REPORT_"`
	} `yaml:"OPA" env-prefix:"OPA_"`
	Reviewer struct {
		Locker struct {
			TableName string `yaml:"TableName" env:"TABLE_NAME"`
		} `yaml:"Locker" env-prefix:"LOCKER_"`
	} `yaml:"Reviewer" env-prefix:"REVIEWER_"`
	GHE struct {
		ServiceAccount string `yaml:"ServiceAccount" env:"SERVICE_ACCOUNT"`
		Hostname       string `yaml:"Hostname" env:"HOSTNAME" env-default:"github.com"`
	} `yaml:"GHE" env-prefix:"GHE_"`
	ConfigStore struct {
		Table   string        `yaml:"Table" env:"TABLE"`
		Refresh time.Duration `yaml:"Refresh" env:"REFRESH"`
	} `yaml:"ConfigStore" env-prefix:"CONFIG_STORE_"`
	Throttler struct {
		Table        string `yaml:"Table" env:"TABLE"`
		PartitionKey string `yaml:"PartitionKey" env:"PARTITION_KEY"`
		SortKey      string `yaml:"SortKey" env:"SORT_KEY"`
		TTLFieldName string `yaml:"TTLFieldName" env:"TTL_FIELD_NAME"`
	} `yaml:"Throttler" env-prefix:"THROTTLER_"`
	Datastore struct {
		Table string        `yaml:"Table" env:"TABLE"`
		TTL   time.Duration `yaml:"TTL" env:"TTL"`
	} `yaml:"Datastore" env-prefix:"DATASTORE_"`
}

func ParseConfigFiles() (*Config, error) {
	var cfg Config

	file := ProcessCLIArgs(&cfg)

	err := cleanenv.ReadConfig(file, &cfg)
	if err != nil {
		log.Err(err).Str("configfile", file).Msg("Error reading configuration from file")
		return nil, err
	}

	return &cfg, nil
}

// ProcessArgs processes and handles CLI arguments
func ProcessCLIArgs(cfg *Config) string {
	var path string

	f := flag.NewFlagSet("pr-bot", flag.ExitOnError)
	f.StringVar(&path, "config", "config.yaml", "file path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		//nolint:all
		fmt.Println()
		//nolint:all
		fmt.Println(envHelp)
	}
	err := f.Parse(os.Args[1:])

	if len(path) < 1 || err != nil {
		f.Usage()
		os.Exit(1)
	}

	return path
}
