package app

import (
	"context"
	"github.com/rmnkmr/go-common/log"
	utils "github.com/rmnkmr/lsp/errors"
	lsp "github.com/rmnkmr/lsp/loanprovider"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) UpdateKycDetails(ctx context.Context, request *lspPb.UploadKycDetailsRequest) (*lspPb.UploadKycDetailsResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "UpdateKycDetails: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "UpdateKycDetails: GetLoanProviderBySlug")
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

	res, err := client.UpdateKycDetails(ctx, request)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "InitiateLoan: error from provider")
		return nil, err
	}

	return res, nil
}
