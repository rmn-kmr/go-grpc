package stashfin

import (
	"context"
	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	api "github.com/rmnkmr/lsp/proto"
)

func (lp ApiClient) CreateLoan(ctx context.Context, request *api.CreateLoanRequest) (*api.CreateLoanResponse, error) {
	applicationID := request.UserLoanDetails.GetLeadId()

	updateBankDetailRes, err := updateBankDetails(ctx, lp, applicationID, request.UserBankDetails)
	if err != nil {
		log.Error(ctx, err, "CreateLoan: updateBankDetails failed", "err", err)
		return nil, err
	}
	log.Info(ctx, "CreateLoan: updateBankDetails res", "updateBankDetailRes", updateBankDetailRes.Message)

	loanStatusReq := &api.CheckLoanStatusRequest{NbfcLeadId: applicationID}
	statusRes, err := lp.CheckLoanStatus(ctx, loanStatusReq)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: CheckLoanStatus failed")
		return nil, err
	}
	createLoanResponse := &api.CreateLoanResponse{}

	if statusRes.Status == api.CheckLoanStatusResponse_APPROVED || statusRes.Status == api.CheckLoanStatusResponse_ELIGIBLE {
		createLoanResponse = &api.CreateLoanResponse{
			NbfcLeadId:      request.UserLoanDetails.LeadId,
			Status:          "SUCCESS",
			NbfcBorrowerId:  request.UserPersonalDetails.NbfcBorrowerId,
			RejectionReason: map[string]string{},
			IsSuccess:       true,
		}
	} else if statusRes.Status == api.CheckLoanStatusResponse_FAILED {
		createLoanResponse = &api.CreateLoanResponse{
			NbfcLeadId:      request.UserLoanDetails.LeadId,
			Status:          "REJECTED",
			NbfcBorrowerId:  request.UserPersonalDetails.NbfcBorrowerId,
			RejectionReason: map[string]string{"status": "Rejected"},
			IsSuccess:       false,
		}
	} else if statusRes.Status != api.CheckLoanStatusResponse_APPROVED {
		customError := errors.New(statusRes.Status.String())
		log.Error(ctx, customError, "CheckLoanLimitStatus: CheckLoanStatus status failed")
		return nil, customError
	}

	return createLoanResponse, nil
}
