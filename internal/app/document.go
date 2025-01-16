package app

import (
	"context"

	utils "github.com/rmnkmr/lsp/errors"
	lsp "github.com/rmnkmr/lsp/loanprovider"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) UploadDocuments(ctx context.Context, request *lspPb.UploadDocumentRequest) (*lspPb.UploadDocumentResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "CreateLoan: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "CreateLoan: GetLoanProviderBySlug")
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

	// amount, _ := strconv.Atoi(req.GetAmount()) // this should be in int64 - instead of typecasting

	_, err = client.UploadDocuments(ctx, request)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "InitiateLoan: error from provider")
		return nil, err
	}

	// todo
	return &lspPb.UploadDocumentResponse{}, nil
}
