package app

import (
	"context"

	lsp "github.com/rmnkmr/lsp/loanprovider"

	utils "github.com/rmnkmr/lsp/errors"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

// get LoanProviders api
func (s *App) InitiateLoan(ctx context.Context, req *lspPb.InitiateLoanRequest) (*lspPb.InitiateLoanResponse, error) {

	provider := req.GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "InitiateLoan: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "InitiateLoan: GetLoanProviderBySlug")
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

	initiateTransactionRequest := &lspPb.InitiateLoanRequest{}

	response, err := client.InitiateLoan(ctx, initiateTransactionRequest)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "InitiateLoan: error from provider")
		return nil, err
	}

	return &lspPb.InitiateLoanResponse{
		LoanId: response.LoanId,
	}, nil
}

func (s *App) CreateLoan(ctx context.Context, request *lspPb.CreateLoanRequest) (*lspPb.CreateLoanResponse, error) {
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
		GCSClient:                     s.GCSClient,
	}

	// amount, _ := strconv.Atoi(req.GetAmount()) // this should be in int64 - instead of typecasting

	res, err := client.CreateLoan(ctx, request)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "CreateLoan: CreateLoan error")
		return nil, err
	}

	// todo
	return res, nil
}

func (s *App) GetLoan(ctx context.Context, request *lspPb.GetLoanRequest) (*lspPb.GetLoanResponse, error) {
	provider := request.GetMetadata().GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "GetLoan: Invalid request")
		return nil, utils.ErrInvalidRequest
	}

	paymentProvider, err := s.Queries.GetLoanProviderBySlug(ctx, provider)
	if err != nil {
		log.Error(ctx, err, "GetLoan: GetLoanProviderBySlug")
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

	getLoanRequest := &lspPb.GetLoanRequest{}

	response, err := client.GetLoan(ctx, getLoanRequest)
	if err != nil && err != utils.ErrProviderErr {
		log.Error(ctx, err, "GetLoan: error from provider")
		return nil, err
	}

	// todo
	return &lspPb.GetLoanResponse{
		Status: response.GetStatus(),
	}, nil
}
