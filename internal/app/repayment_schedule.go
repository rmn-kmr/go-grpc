package app

import (
	"context"

	utils "github.com/rmnkmr/lsp/errors"
	lsp "github.com/rmnkmr/lsp/loanprovider"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

func (s *App) CreateRepaymentSchedule(ctx context.Context, request *lspPb.CreateRepaymentScheduleRequest) (*lspPb.CreateRepaymentScheduleResponse, error) {
	provider := request.GetProvider().String()

	if provider == "" {
		log.Error(ctx, utils.ErrInvalidRequest, "CreateLead: Invalid request")
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
	client.CreateRepaymentSchedule(ctx, request)
	return nil, nil
}
