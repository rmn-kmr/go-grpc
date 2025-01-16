package app

import (
	"cloud.google.com/go/bigquery"
	gstorage "cloud.google.com/go/storage"
	"context"
	lspConfig "github.com/rmnkmr/lsp/config"
	lspDB "github.com/rmnkmr/lsp/internal/db"
	lspApi "github.com/rmnkmr/lsp/proto"
	pubsub "github.com/rmnkmr/lsp/pubsub"
	protonium "github.com/rmnkmr/protonium/http"
)

func New(s *App) protonium.Application {
	return &handler{
		api: s,
	}
}

type handler struct {
	api           lspApi.APIServer
	pathParamFunc protonium.PathParamFunc
}

func (h *handler) Initialize(fn protonium.PathParamFunc) {
	h.pathParamFunc = fn
}

func (h *handler) Routes() []*protonium.Route {
	panic("implement me")
}

type App struct {
	Queries     *lspDB.Queries
	Publisher   pubsub.Publisher
	Environment string
	Config      lspConfig.ServiceConfig
	BqClient    *bigquery.Client
	GCSClient   gstorage.Client
}

func (s *App) BureauPull(ctx context.Context, request *lspApi.BureauPullRequest) (*lspApi.BureauPullResponse, error) {
	//TODO implement me
	panic("implement me")
}
