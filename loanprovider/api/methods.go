package api

import (
	"context"
	lsp "github.com/rmnkmr/lsp/proto"
)

type LoanProviderAPI interface {
	GetTokenFromProvider(context.Context) (*AuthToken, error)
	OnCallBack(ctx context.Context, request *OnCallBackRequest) (*CallbackResponse, error)
	InitiateLoan(ctx context.Context, request *lsp.InitiateLoanRequest) (*lsp.InitiateLoanResponse, error)
	CreateLead(ctx context.Context, request *lsp.CreateLeadRequest) (*lsp.CreateLeadResponse, error)
	SaveRepayment(ctx context.Context, request *lsp.SaveRepaymentRequest) (*lsp.SaveRepaymentResponse, error)
	CreateLoan(ctx context.Context, request *lsp.CreateLoanRequest) (*lsp.CreateLoanResponse, error)
	UploadDocuments(ctx context.Context, request *lsp.UploadDocumentRequest) (*lsp.UploadDocumentResponse, error)
	GetLoan(ctx context.Context, request *lsp.GetLoanRequest) (*lsp.GetLoanResponse, error)
	CreateRepaymentSchedule(context.Context, *lsp.CreateRepaymentScheduleRequest) (*lsp.CreateRepaymentScheduleResponse, error)
	CheckLoanLimitStatus(ctx context.Context, request *lsp.LoanLimitStatusRequest) (*lsp.LoanLimitStatusResponse, error)
	DisburseLoan(ctx context.Context, request *lsp.DisbursalLoanRequest) (*lsp.DisbursalLoanResponse, error)
	CheckLoanStatus(ctx context.Context, request *lsp.CheckLoanStatusRequest) (*lsp.CheckLoanStatusResponse, error)
	PaymentUrl(ctx context.Context, request *lsp.PaymentUrlRequest) (*lsp.PaymentUrlResponse, error)
	UpdateKycDetails(ctx context.Context, request *lsp.UploadKycDetailsRequest) (*lsp.UploadKycDetailsResponse, error)
}
