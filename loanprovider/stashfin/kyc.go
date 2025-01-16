package stashfin

import (
	"context"
	"github.com/rmnkmr/go-common/log"
	api "github.com/rmnkmr/lsp/proto"
)

func (lp ApiClient) UpdateKycDetails(ctx context.Context, request *api.UploadKycDetailsRequest) (*api.UploadKycDetailsResponse, error) {
	applicationId := request.LoanId
	updateAddressRes, err := updateAddress(ctx, lp, applicationId, request.UserPersonalDetails)
	if err != nil {
		log.Error(ctx, err, "UpdateKycDetails: updateAddress failed", "err", err)
		return nil, err
	}
	log.Info(ctx, "UpdateKycDetails: updateAddress res", "updateAddressRes", updateAddressRes.Message)

	_, err = uploadKYCDocuments(ctx, lp, applicationId, request.UserKycDetails, request.UserPersonalDetails.PanMetadata)
	if err != nil {
		log.Error(ctx, err, "UpdateKycDetails: UploadDocuments failed", "err", err)
		return nil, err
	}
	log.Info(ctx, "UpdateKycDetails: UploadDocuments success")

	// Update id details
	_, err = updateIdDetail(ctx, lp, applicationId, request.UserPersonalDetails)
	if err != nil {
		log.Error(ctx, err, "UpdateKycDetails: updateIdDetail failed", "err", err)
		return nil, err
	}
	log.Info(ctx, "UpdateKycDetails: updateIdDetail success")
	return &api.UploadKycDetailsResponse{Status: true}, nil
}
