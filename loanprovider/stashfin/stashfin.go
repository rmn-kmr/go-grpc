package stashfin

import (
	gstorage "cloud.google.com/go/storage"
	"context"
	api2 "github.com/rmnkmr/lsp/loanprovider/api"
	api "github.com/rmnkmr/lsp/proto"
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

func (lp ApiClient) UploadDocuments(ctx context.Context, request *api.UploadDocumentRequest) (*api.UploadDocumentResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ api2.LoanProviderAPI = &ApiClient{}

func (lp ApiClient) SaveRepayment(ctx context.Context, request *api.SaveRepaymentRequest) (*api.SaveRepaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (lp ApiClient) GetLoan(ctx context.Context, request *api.GetLoanRequest) (*api.GetLoanResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (lp ApiClient) CreateRepaymentSchedule(ctx context.Context, request *api.CreateRepaymentScheduleRequest) (*api.CreateRepaymentScheduleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (lp ApiClient) GetTokenFromProvider(ctx context.Context) (*api2.AuthToken, error) {
	return lp.login(ctx)
}

func (lp ApiClient) OnCallBack(ctx context.Context, request *api2.OnCallBackRequest) (*api2.CallbackResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (lp ApiClient) InitiateLoan(ctx context.Context, request *api.InitiateLoanRequest) (*api.InitiateLoanResponse, error) {
	//TODO implement me
	panic("implement me")
}
