package loanprovider

import (
	"context"

	gstorage "cloud.google.com/go/storage"
	"github.com/rmnkmr/lsp/loanprovider/stashfin"
	lsp "github.com/rmnkmr/lsp/proto"

	api "github.com/rmnkmr/lsp/loanprovider/api"
	log "github.com/rmnkmr/lsp/log"
)

type ApiClient struct {
	Environment                   string
	Provider                      string
	ProviderID                    string
	ProviderAccountId             string
	AccountId                     string
	ProviderBaseUrl               string
	PreExpirationTokenRefreshMins int32
	Token                         string
	ClientId                      string
	ClientSecret                  string
	GCSClient                     gstorage.Client
}

var _ api.LoanProviderAPI = &ApiClient{}

func (lp ApiClient) GetProviderClient(ctx context.Context) api.LoanProviderAPI {
	switch lp.Provider {
	case lsp.Provider_STASHFIN.String():
		return stashfin.ApiClient(lp)
	default:
		return nil
	}
}

// get token from provider api
func (lp ApiClient) GetTokenFromProvider(ctx context.Context) (*api.AuthToken, error) {
	// implement me
	log.Debug(ctx, "GetTokenFromProvider: implement me")
	return &api.AuthToken{}, nil
}

func (lp ApiClient) CreateLead(ctx context.Context, request *lsp.CreateLeadRequest) (*lsp.CreateLeadResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.CreateLead(ctx, request)
}

func (lp ApiClient) InitiateLoan(ctx context.Context, request *lsp.InitiateLoanRequest) (*lsp.InitiateLoanResponse, error) {
	//TODO implement me
	client := lp.GetProviderClient(ctx)
	return client.InitiateLoan(ctx, request)
}

func (lp ApiClient) SaveRepayment(ctx context.Context, request *lsp.SaveRepaymentRequest) (*lsp.SaveRepaymentResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.SaveRepayment(ctx, request)
}

func (lp ApiClient) OnCallBack(ctx context.Context, request *api.OnCallBackRequest) (*api.CallbackResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.OnCallBack(ctx, request)
}

func (lp ApiClient) UploadDocuments(ctx context.Context, request *lsp.UploadDocumentRequest) (*lsp.UploadDocumentResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.UploadDocuments(ctx, request)
}

func (lp ApiClient) GetLoan(ctx context.Context, request *lsp.GetLoanRequest) (*lsp.GetLoanResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.GetLoan(ctx, request)
}

func (lp ApiClient) CreateLoan(ctx context.Context, request *lsp.CreateLoanRequest) (*lsp.CreateLoanResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.CreateLoan(ctx, request)
}

func (lp ApiClient) CreateRepaymentSchedule(ctx context.Context, request *lsp.CreateRepaymentScheduleRequest) (*lsp.CreateRepaymentScheduleResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.CreateRepaymentSchedule(ctx, request)
}
func (lp ApiClient) CheckLoanLimitStatus(ctx context.Context, request *lsp.LoanLimitStatusRequest) (*lsp.LoanLimitStatusResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.CheckLoanLimitStatus(ctx, request)
}
func (lp ApiClient) DisburseLoan(ctx context.Context, request *lsp.DisbursalLoanRequest) (*lsp.DisbursalLoanResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.DisburseLoan(ctx, request)
}

func (lp ApiClient) CheckLoanStatus(ctx context.Context, request *lsp.CheckLoanStatusRequest) (*lsp.CheckLoanStatusResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.CheckLoanStatus(ctx, request)
}

func (lp ApiClient) PaymentUrl(ctx context.Context, request *lsp.PaymentUrlRequest) (*lsp.PaymentUrlResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.PaymentUrl(ctx, request)
}

func (lp ApiClient) UpdateKycDetails(ctx context.Context, request *lsp.UploadKycDetailsRequest) (*lsp.UploadKycDetailsResponse, error) {
	client := lp.GetProviderClient(ctx)
	return client.UpdateKycDetails(ctx, request)
}
