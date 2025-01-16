package app

import (
	"context"

	lsp "github.com/rmnkmr/lsp/loanprovider"

	utils "github.com/rmnkmr/lsp/errors"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) CheckLoanStatus(ctx context.Context, request *lspPb.CheckLoanStatusRequest) (*lspPb.CheckLoanStatusResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "CheckLoanStatus: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "CheckLoanStatus: GetLoanProviderBySlug")
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
		GCSClient:                     s.GCSClient,
	}

	res, err := client.CheckLoanStatus(ctx, request)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "CheckLoanStatus: CheckLoanStatus error")
		return nil, err
	}
	return res, nil
}
