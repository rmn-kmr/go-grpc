package app

import (
	"context"
	"github.com/rmnkmr/lsp/log"
	lspPb "github.com/rmnkmr/lsp/proto"
)

// get LoanProviders api
func (s *App) LoanProviders(ctx context.Context, req *lspPb.LoanProvidersRequest) (*lspPb.LoanProvidersResponse, error) {
	providers, err := s.Queries.GetLoanProviders(ctx)
	if err != nil {
		return nil, err
	}
	var paymentProviders []*lspPb.LoanProviderResponse
	for _, p := range providers {
		paymentProvider := &lspPb.LoanProviderResponse{
			Provider:       p.Provider,
			ProviderId:     p.ProviderID,
			ProviderName:   p.ProviderName,
			ProviderType:   p.ProviderType,
			ProviderConfig: &lspPb.ProviderConfig{},
		}
		paymentProviders = append(paymentProviders, paymentProvider)
		log.Debug(ctx, "LoanProviders", "provider", p)
	}
	return &lspPb.LoanProvidersResponse{
		Providers: paymentProviders,
	}, nil
}
