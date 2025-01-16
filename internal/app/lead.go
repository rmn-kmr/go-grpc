package app

import (
	"context"

	utils "github.com/rmnkmr/lsp/errors"
	lsp "github.com/rmnkmr/lsp/loanprovider"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) CreateLead(ctx context.Context, request *lspPb.CreateLeadRequest) (*lspPb.CreateLeadResponse, error) {
	provider := request.GetMetadata().GetProvider().String()
	log.Info(ctx, "CreateLead: ", "provider", provider, "request", request)
	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "CreateLead: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "CreateLead: GetLoanProviderBySlug")
		return nil, err
	}

	client := lsp.ApiClient{
		Provider:                      paymentProvider.Provider,
		ProviderID:                    paymentProvider.ProviderID,
		ProviderBaseUrl:               paymentProvider.ApiBaseUrl,
		ClientId:                      paymentProvider.ApiKey,
		ClientSecret:                  paymentProvider.ApiSecret,
		Environment:                   s.Environment,
		PreExpirationTokenRefreshMins: paymentProvider.PreExpirationTokenRefreshMins,
		AccountId:                     paymentProvider.ProviderAccountID,
	}
	res, err := client.CreateLead(ctx, request)
	if err != nil {
		log.Error(ctx, err, "InitiateLoan: CreateLead")
		return nil, err
	}
	log.Info(ctx, "CreateLead: ", "provider", provider, "res", res)
	return res, err
}

func (s *App) CheckLoanLimitStatus(ctx context.Context, request *lspPb.LoanLimitStatusRequest) (*lspPb.LoanLimitStatusResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "CheckLoanLimitStatus: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: GetLoanProviderBySlug")
		return nil, err
	}

	client := lsp.ApiClient{
		Provider:                      paymentProvider.Provider,
		ProviderID:                    paymentProvider.ProviderID,
		ProviderBaseUrl:               paymentProvider.ApiBaseUrl,
		ClientId:                      paymentProvider.ApiKey,
		ClientSecret:                  paymentProvider.ApiSecret,
		Environment:                   s.Environment,
		PreExpirationTokenRefreshMins: paymentProvider.PreExpirationTokenRefreshMins,
		AccountId:                     paymentProvider.ProviderAccountID,
	}
	res, err := client.CheckLoanLimitStatus(ctx, request)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: CheckLoanLimitStatus")
		return nil, err
	}
	return res, err
}
