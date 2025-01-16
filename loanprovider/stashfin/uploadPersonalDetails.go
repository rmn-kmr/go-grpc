package stashfin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	api "github.com/rmnkmr/lsp/proto"
)

const (
	update_address              = "/v3/update-address"
	update_personal_information = "/v3/update-personal-information"
	update_professional_details = "/v3/update-professional-details"
	update_bank_details         = "/v3/update-bank-details"
	update_id_details           = "/v3/update-id-details"
)

type updateAddressPayload struct {
	ApplicationId  interface{}     `json:"application_id"`
	AddressDetails []AddressDetail `json:"address_details"`
}

type AddressDetail struct {
	AddressType   int32  `json:"address_type"`
	AddressLine1  string `json:"address_line_1"`
	Pincode       int32  `json:"pincode"`
	AddressLine2  string `json:"address_line_2"`
	OwnershipType int32  `json:"ownership_type"`
}

type updatePersonalDetailResponse struct {
	Status  bool   `json:"status"`
	Results string `json:"results"`
}

type updatePersonalInfoPayload struct {
	ApplicationId string `json:"application_id"`
	FatherName    string `json:"father_name"`
	MaritalStatus int    `json:"marital_status"`
}

type updateProfessionalInfoPayload struct {
	ApplicationId string `json:"application_id"`
	CompanyName   string `json:"company_name"`
	PartnerScore  int64  `json:"partner_score"`
}

type updateBankDetailsPayload struct {
	ApplicationId string `json:"application_id"`
	IfscCode      string `json:"ifsc_code"`
	AccountNumber string `json:"account_number"`
}

type updateIDDetailsPayload struct {
	ApplicationId string     `json:"application_id"`
	IdDetails     []IdDetail `json:"id_details"`
}

type updateIdDetailResponse struct {
	Status  bool   `json:"status"`
	Results string `json:"results"`
}

type IdDetail struct {
	IdType       int    `json:"id_type"`
	IdNumber     string `json:"id_number"`
	DocumentName string `json:"document_name"`
}

func updateAddress(ctx context.Context, lp ApiClient, applicationId string, personalDetail *api.UserPersonalDetails) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, update_address)
	method := http.MethodPost
	var err error
	var pincode int
	var address string
	if personalDetail.Pincode != "" {
		pincode, err = strconv.Atoi(personalDetail.Pincode)
		if err != nil {
			log.Error(ctx, err, "updateAddress: Pincode int type conversion failed", "err", err, "Pincode", personalDetail.Pincode)
			return nil, err
		}
	}

	if personalDetail.Address != "" {
		address = personalDetail.Address
	}
	address1 := address
	address2 := ""
	if len(address1) > 40 {
		address1 = address[:40]
		address2 = address[40:]
	}
	if len(address2) > 40 {
		address2 = address2[:40]
	}
	payload := &updateAddressPayload{
		ApplicationId: applicationId,
		AddressDetails: []AddressDetail{
			{
				AddressType:   Permanent_Address,
				AddressLine1:  address1,
				Pincode:       int32(pincode),
				AddressLine2:  address2,
				OwnershipType: Current_Residence_Address,
			},
		},
	}
	apiResponse := &updatePersonalDetailResponse{}
	err = httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "updateAddress: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	if apiResponse == nil || !apiResponse.Status {
		log.Info(ctx, "updateAddress: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, errors.New("updateAddress: API call failed")
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: apiResponse.Results,
	}, nil
}

//func updatePersonalInformation(ctx context.Context, lp ApiClient, applicationId string, personalDetail *api.UserPersonalDetails) (*api.APISuccessResponse, error) {
//	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, update_personal_information)
//	method := "POST"
//	fatherName := "fatherName"
//	if personalDetail.FathersName != "" {
//		fatherName = personalDetail.FathersName
//	}
//	payload := &updatePersonalInfoPayload{
//		ApplicationId: applicationId,
//		FatherName:    fatherName,
//		MaritalStatus: Marital_Status_single,
//	}
//	apiResponse := &updatePersonalDetailResponse{}
//	err := httpCall(ctx, lp, url, method, payload, apiResponse)
//	if err != nil {
//		log.Error(ctx, err, "updatePersonalInformation: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
//		return nil, err
//	}
//	return &api.APISuccessResponse{
//		Status:  apiResponse.Status,
//		Message: apiResponse.Results,
//	}, nil
//}

func updateProfessionalDetails(ctx context.Context, lp ApiClient, req *updateProfessionalInfoPayload) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, update_professional_details)
	method := http.MethodPost

	payload := &updateProfessionalInfoPayload{
		ApplicationId: req.ApplicationId,
		CompanyName:   Self_Employed,
		PartnerScore:  req.PartnerScore,
	}

	apiResponse := &updatePersonalDetailResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "updateProfessionalDetails: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	if apiResponse == nil || apiResponse.Status == false {
		log.Info(ctx, "updateProfessionalDetails: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, errors.New("updateProfessionalDetails: API call failed")
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: apiResponse.Results,
	}, nil
}

func updateBankDetails(ctx context.Context, lp ApiClient, applicationId string, bankDetail *api.UserBankDetails) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, update_bank_details)
	method := http.MethodPost
	var ifsc, accNo string
	if bankDetail.BankIfsc != "" {
		ifsc = bankDetail.BankIfsc
	}

	if bankDetail.BankAccNum != "" {
		accNo = bankDetail.BankAccNum
	}
	payload := &updateBankDetailsPayload{
		ApplicationId: applicationId,
		IfscCode:      ifsc,
		AccountNumber: accNo,
	}
	apiResponse := &updatePersonalDetailResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "updateBankDetails: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	if apiResponse == nil || apiResponse.Status == false {
		log.Info(ctx, "updateBankDetails: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, errors.New("updateBankDetails: API call failed")
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: apiResponse.Results,
	}, nil
}

func updateIdDetail(ctx context.Context, lp ApiClient, applicationId string,
	userPersonalDetail *api.UserPersonalDetails) (*api.APISuccessResponse, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, update_id_details)
	method := http.MethodPost

	payload := &updateIDDetailsPayload{
		ApplicationId: applicationId,
		IdDetails: []IdDetail{{
			DocumentName: string(updateDocumentTyoe_aadhaar),
			IdNumber:     fmt.Sprintf("XXXXXXXX%s", userPersonalDetail.AddrProofNumber[8:12]),
			IdType:       int(updateDocumentID_aadhaar),
		}},
	}
	apiResponse := &updateIdDetailResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "updateIdDetail: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, err
	}
	if apiResponse == nil || apiResponse.Status == false {
		log.Info(ctx, "updateIdDetail: API call failed", "err", err, "url", url, "apiResponse", apiResponse)
		return nil, errors.New("updateIdDetail: API call failed")
	}
	return &api.APISuccessResponse{
		Status:  apiResponse.Status,
		Message: apiResponse.Results,
	}, nil
}
