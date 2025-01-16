package stashfin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rmnkmr/lsp/log"
	api "github.com/rmnkmr/lsp/proto"
	"github.com/rmnkmr/lsp/utils"
)

const (
	save_amortization      = "/v3/save-amortization"
	save_and_approve_amort = "/v3/save-and-approve-amort"
	payment_url            = "/v3/get-payment-url"
)

type SaveAmortizationPayload struct {
	DisbursalDate        string        `json:"disbursal_date"`
	ApprovedAmount       float32       `json:"approved_amount"`
	Tenure               int           `json:"tenure"`
	FirstEmiDate         string        `json:"first_emi_date"`
	ProcessingFee        float32       `json:"processing_fee"`
	Roi                  float32       `json:"roi"`
	OtherFees            int           `json:"other_fees"`
	OtherFeesBifurcation []interface{} `json:"other_fees_bifurcation"`
	LeadId               string        `json:"lead_id"`
}

type saveAndApproveAmortPayload struct {
	LeadId    string       `json:"lead_id"`
	Finalamot []AmortEntry `json:"finalamot"`
}

type AmortEntry struct {
	Principal       float32 `json:"principal"`
	Discount        float32 `json:"discount"`
	Date            string  `json:"date"`
	Interest        float32 `json:"interest"`
	StartingBalance float32 `json:"starting_balance"`
	ClosingBalance  float32 `json:"closing_balance"`
	Amount          float32 `json:"amount"`
}

type SaveAmortResponse struct {
	Status  bool `json:"status"`
	Results struct {
		UpfrontInterest          int     `json:"upfront_interest"`
		Roi                      float64 `json:"roi"`
		DisbursalDate            string  `json:"disbursal_date"`
		Tenure                   int     `json:"tenure"`
		ProcessingFee            int     `json:"processing_fee"`
		FirstEmiDate             string  `json:"first_emi_date"`
		ApprovedAmount           int     `json:"approved_amount"`
		LenderDisbursedAmount    int     `json:"lender_disbursed_amount"`
		LenderRoi                float64 `json:"lender_roi"`
		GstCharge                int     `json:"gst_charge"`
		LenderDisbursedAmountGst int     `json:"lender_disbursed_amount_gst"`
		LenderSPDC               int     `json:"lender_SPDC"`
		Schedule                 []struct {
			Principal       int    `json:"principal"`
			Discount        int    `json:"discount"`
			Date            string `json:"date"`
			Interest        int    `json:"interest"`
			StartingBalance int    `json:"starting_balance"`
			ClosingBalance  int    `json:"closing_balance"`
			Amount          int    `json:"amount"`
		} `json:"schedule"`
	} `json:"results"`
}

type saveAndApproveAmortResponse struct {
	Status  bool `json:"status"`
	Results struct {
		Success string `json:"success"`
	} `json:"results"`
}

type paymentUrlRequest struct {
	ApplicationId string  `json:"application_id"`
	Amount        float32 `json:"amount"`
}

type paymentUrlResponse struct {
	Status  bool   `json:"status"`
	Results string `json:"results"`
}

func saveAmortization(ctx context.Context, lp ApiClient, applicationId string, userLoanDetail *api.UserLoanDetails) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, save_amortization)
	method := http.MethodPost
	disbursalDate := utils.GetTimeFromEpoch(userLoanDetail.DisbursalDate).Format(utils.YYYY_MM_DD)
	firstEmiDate := utils.GetTimeFromEpoch(userLoanDetail.FirstEmiDate).Format(utils.YYYY_MM_DD)

	roundedMonths := CalculateMonthsForNumberOfDays(userLoanDetail.Installments)
	roundedMonthlyInterest := CalculateMonthlyInterest(userLoanDetail.LoanInterestPercentage)

	payload := &SaveAmortizationPayload{
		DisbursalDate:  disbursalDate,
		ApprovedAmount: userLoanDetail.LoanAmount,
		Tenure:         roundedMonths,
		FirstEmiDate:   firstEmiDate,
		ProcessingFee:  userLoanDetail.ProcessingFeePercentage,
		Roi:            roundedMonthlyInterest,
		OtherFees:      0,
		LeadId:         applicationId,
	}
	apiResponse := &SaveAmortResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "saveAmortization: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: "Success",
	}, nil
}

func saveAndApproveAmortization(ctx context.Context, lp ApiClient, applicationId string, userLoanDetail *api.UserLoanDetails) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, save_and_approve_amort)
	method := http.MethodPost
	var amortRecords []AmortEntry
	startingBalance := userLoanDetail.LoanAmount
	for _, record := range userLoanDetail.InstallmentData {
		closingBalance := startingBalance - record.PrincipalDue
		newEntry := AmortEntry{
			Principal:       record.PrincipalDue,
			Discount:        0,
			Date:            utils.GetTimeFromEpoch(record.InstallmentDate).Format(utils.YYYY_MM_DD),
			Interest:        record.InterestDue,
			StartingBalance: startingBalance,
			ClosingBalance:  closingBalance,
			Amount:          record.PrincipalDue + record.InterestDue,
		}
		amortRecords = append(amortRecords, newEntry)
		startingBalance = closingBalance
	}

	payload := saveAndApproveAmortPayload{
		LeadId:    applicationId,
		Finalamot: amortRecords,
	}
	apiResponse := &saveAndApproveAmortResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "saveAndApproveAmortization: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: apiResponse.Results.Success,
	}, nil
}

func (lp ApiClient) PaymentUrl(ctx context.Context, request *api.PaymentUrlRequest) (*api.PaymentUrlResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, payment_url)
	method := "POST"
	payload := paymentUrlRequest{
		ApplicationId: request.NbfcLeadId,
		Amount:        request.Amount,
	}
	apiResponse := &paymentUrlResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "PaymentUrl: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	return &api.PaymentUrlResponse{
		RepaymentUrl: apiResponse.Results,
		Status:       apiResponse.Status,
	}, nil
}
