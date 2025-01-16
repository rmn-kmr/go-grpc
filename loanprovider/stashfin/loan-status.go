package stashfin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	api "github.com/rmnkmr/lsp/proto"
)

type postApprovalStatusResponse struct {
	Status  bool `json:"status"`
	Results struct {
		IsSelfieUploaded                    bool `json:"is_selfie_uploaded"`
		Okyc                                bool `json:"okyc"`
		ApprovedScreenAndAcceptancePage     bool `json:"approved_screen_and_acceptance_page"`
		Enach                               bool `json:"enach"`
		AddressVerificationForCardDelievery bool `json:"address_verification_for_card_delievery"`
	} `json:"results"`
}

type checkLimitStatusResponse struct {
	Status  bool        `json:"status"`
	Results string      `json:"results"`
	Errors  interface{} `json:"errors"`
}

type LimitStatusResult struct {
	Status interface{} `json:"status"`
	Amount int64       `json:"approvedAmount"`
	Msg    string      `json:"msg"`
}

func (lp ApiClient) CheckLoanLimitStatus(ctx context.Context, request *api.LoanLimitStatusRequest) (*api.LoanLimitStatusResponse, error) {
	payload := &updateProfessionalInfoPayload{
		ApplicationId: request.NbfcLoanId,
		CompanyName:   request.CompanyName,
		PartnerScore:  request.PartnerScore,
	}
	res, err := updateProfessionalDetails(ctx, lp, payload)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: updateProfessionalDetails failed", "req", payload, "res", res)
		return nil, err
	}

	url := fmt.Sprintf("%s%s%s", lp.ProviderBaseUrl, check_limit_status, request.NbfcLoanId)
	method := http.MethodGet
	reqPayload := EmptyRequest{}
	response := &checkLimitStatusResponse{}
	err = httpCall(ctx, lp, url, method, reqPayload, response)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: httpCall failed", "err", err, "url", url, "payload", payload)
		return nil, err
	}

	if response == nil || !response.Status || response.Results == "" {
		log.Info(ctx, "CheckLoanLimitStatus: lead creation failed", "response", response, "url", url, "payload", payload)
		return nil, errors.New("CheckLoanLimitStatus: API call failed")
	}

	var result LimitStatusResult
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(response.Results), &result)
	if err != nil {
		log.Error(ctx, err, "CheckLoanLimitStatus: json Unmarshal failed")
		return nil, err
	}

	llsRes := &api.LoanLimitStatusResponse{}

	if result.Status == PASSED {
		llsRes = &api.LoanLimitStatusResponse{
			IsSuccess:  response.Status,
			IsApproved: result.Amount > 0,
			MaxAmount:  result.Amount,
		}
	} else if result.Status == REJECTED {
		llsRes = &api.LoanLimitStatusResponse{
			IsSuccess:  false,
			IsApproved: false,
			IsRejected: true,
		}
	} else {
		llsRes = &api.LoanLimitStatusResponse{
			IsSuccess:  false,
			IsApproved: false,
			IsRejected: false,
		}
	}

	if !llsRes.IsApproved {
		llsRes.Msg = response.Results
	}
	return llsRes, nil
}

func (lp ApiClient) CheckLoanStatus(ctx context.Context, req *api.CheckLoanStatusRequest) (*api.CheckLoanStatusResponse, error) {
	url := fmt.Sprintf("%s%s%s", lp.ProviderBaseUrl, check_status, req.NbfcLeadId)
	method := http.MethodGet

	apiResponse := &api.CheckStatusResponse{}
	err := httpCall(ctx, lp, url, method, nil, apiResponse)
	if err != nil {
		log.Error(ctx, err, "checkStatus: API call failed", "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	if apiResponse == nil || !apiResponse.Status {
		log.Info(ctx, "checkStatus: API call failed", "url", url, "apiResponse", apiResponse)
		return nil, errors.New("checkStatus: API call failed")
	}
	loanStatus := &api.CheckLoanStatusResponse{}
	switch apiResponse.Results.ApplicationStatus {
	case "Incomplete Application":
		loanStatus.Status = api.CheckLoanStatusResponse_INCOMPLETE_APPLICATION
	case "Under Processing":
		loanStatus.Status = api.CheckLoanStatusResponse_UNDER_PROCESSING
	case "Approved":
		loanStatus.Status = api.CheckLoanStatusResponse_APPROVED
	case "Disbursed":
		loanStatus.Status = api.CheckLoanStatusResponse_SUCCESS
	case "Eligible":
		loanStatus.Status = api.CheckLoanStatusResponse_ELIGIBLE
	case "Rejected":
		loanStatus.Status = api.CheckLoanStatusResponse_FAILED
	}
	return loanStatus, nil
}
