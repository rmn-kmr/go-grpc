package app

import (
	"context"

	utils "github.com/rmnkmr/lsp/errors"
	lsp "github.com/rmnkmr/lsp/loanprovider"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) DisburseLoan(ctx context.Context, request *lspPb.DisbursalLoanRequest) (*lspPb.DisbursalLoanResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "DisburseLoan: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "DisburseLoan: GetLoanProviderBySlug")
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
	res, err := client.DisburseLoan(ctx, request)
	if err != nil {
		log.Error(ctx, err, "DisburseLoan: DisburseLoan")
		return nil, err
	}
	return res, nil
}
