package entrypoint

import (
	"context"
	"database/sql"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	gstorage "cloud.google.com/go/storage"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rmnkmr/go-common"
	configManager "github.com/rmnkmr/go-common/config"
	"github.com/rmnkmr/go-common/newrelic"
	"github.com/rmnkmr/go-common/postgres"
	"github.com/rmnkmr/lsp/internal/app"
	"github.com/rmnkmr/lsp/internal/callback"
	lspDB "github.com/rmnkmr/lsp/internal/db"
	"github.com/rmnkmr/lsp/internal/grpc"
	"github.com/rmnkmr/lsp/log"
	"github.com/rmnkmr/lsp/pubsub"
	"github.com/rmnkmr/protonium"
)

const (
	Name = "lsp"
)

var conf Configuration
var nrApp *nr.Application
var a *app.App
var bgClient *bigquery.Client

func Execute() {

	// service runtime context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// config
	confManager, err := configManager.Load()
	if err != nil {
		panic(err)
	}

	err = confManager.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	log.InitLogger(conf.Log.Level)

	nrApp = newrelic.App(Name, conf.Env, conf.NewRelic)
	engineConf := conf.Service

	pubsubClient, err := pubsub.NewClient(ctx, conf.Service.Gcp.ProjectID)
	if err != nil {
		panic(err)
	}
	queriesClient := lspDB.New(initSqlDatabase(conf.Env, conf))

	a = &app.App{
		Queries:     queriesClient,
		Environment: conf.Env.String(),
		Publisher:   pubsub.Publisher{Environment: conf.Env.String(), ProjectID: conf.Service.Gcp.ProjectID, Client: pubsubClient},
		Config:      engineConf,
		GCSClient:   initGCSClient(context.Background()),
	}

	c := &callback.Callback{
		App: a,
	}

	bgClient, err = initBigQueryClient(a.Config.Gcp.BigQueryProjectId)
	if err != nil {
		panic(err)
	}

	a.BqClient = bgClient

	s := protonium.New(Name, conf.Env, nrApp,
		append(
			callback.ListenerServer(ctx, c),
			grpc.New(a, conf.Port),
		)...,
	)
	s.Run()
}

func initSqlDatabase(env common.Env, conf Configuration) *sql.DB {
	return postgres.New(conf.Database)
}

func initBigQueryClient(projectID string) (*bigquery.Client, error) {
	ctx := context.Background()
	var err error
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func initGCSClient(ctx context.Context) gstorage.Client {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Error(ctx, err, "initGCSClient: Failed to init GCS Client")
		panic(err)
	}
	log.Info(ctx, "gcs client initiated")
	return *client
}
