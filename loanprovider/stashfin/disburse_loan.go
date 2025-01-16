package stashfin

import (
	"context"

	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	lsp "github.com/rmnkmr/lsp/proto"
)

func (lp ApiClient) DisburseLoan(ctx context.Context, req *lsp.DisbursalLoanRequest) (*lsp.DisbursalLoanResponse, error) {
	// Upload Loan documents
	_, err := uploadLoanDocuments(ctx, lp, req.NbfcLeadId, req.Documents)
	if err != nil {
		log.Error(ctx, err, "CreateLoan: UploadDocuments failed", "err", err)
		return nil, err
	}
	log.Info(ctx, "DisburseLoan: UploadDocuments success")

	// Create loan
	res, err := saveAmortization(ctx, lp, req.NbfcLeadId, req.UserLoanDetail)
	if err != nil {
		log.Error(ctx, err, "DisburseLoan: saveAmortization failed")
		return nil, err
	}
	if !res.Status {
		err := errors.New(res.Message)
		log.Error(ctx, err, "DisburseLoan: saveAmortization API failed", "res", res)
		return nil, err
	}

	// Save amort schedule
	res, err = saveAndApproveAmortization(ctx, lp, req.NbfcLeadId, req.UserLoanDetail)
	if err != nil {
		log.Error(ctx, err, "DisburseLoan: saveAndApproveAmortization failed")
		return nil, err
	}
	if !res.Status {
		err := errors.New(res.Message)
		log.Error(ctx, err, "DisburseLoan: saveAndApproveAmortization API failed", "res", res)
		return nil, err
	}

	return &lsp.DisbursalLoanResponse{Status: true}, err
}
