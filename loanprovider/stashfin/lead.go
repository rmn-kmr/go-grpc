package stashfin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	api "github.com/rmnkmr/lsp/proto"
	"github.com/rmnkmr/lsp/utils"
)

type EmptyRequest struct {
}

type initiateApplicationPayload struct {
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Gender         string `json:"gender"`
	Dob            string `json:"dob"`
	PanNumber      string `json:"pan_number"`
	Pincode        int    `json:"pincode"`
	Income         int    `json:"income"`
	EmploymentType int    `json:"employment_type"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

type initiateApplicationResponse struct {
	Status  bool `json:"status"`
	Results struct {
		ApplicationId string `json:"application_id"`
		CustomerId    int    `json:"customer_id"`
		RedirectUrl   string `json:"redirect_url"`
	} `json:"results"`
	Errors interface{} `json:"errors"`
}

type checkDuplicatePayload struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type checkDuplicateResponse struct {
	Status  bool `json:"status"`
	Results bool `json:"results"`
}

func (lp ApiClient) CreateLead(ctx context.Context, request *api.CreateLeadRequest) (*api.CreateLeadResponse, error) {
	isExistingLead, err := checkDuplicate(ctx, lp, request)
	if err != nil {
		log.Error(ctx, err, "CreateLead: checkDuplicate failed", "err", err)
		return nil, err
	}

	if isExistingLead {
		err = errors.From(400, "lead already exists")
		log.Info(ctx, "CreateLead: checkDuplicate", "isExistingLead", isExistingLead)
		return nil, err
	}

	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, initiate_application)
	method := "POST"

	payload, err := lp.toCreateLeadPayload(request)
	if err != nil {
		log.Error(ctx, err, "CreateLead: toCreateLeadPayload failed", "err", err)
		return nil, err
	}
	response := &initiateApplicationResponse{}
	err = httpCall(ctx, lp, url, method, payload, response)
	if err != nil {
		log.Error(ctx, err, "CreateLead: httpCall failed", "err", err, "url", url, "payload", payload)
		return nil, err
	}

	if response == nil || !response.Status {
		log.Info(ctx, "CreateLead: lead creation failed", "response", response, "url", url, "payload", payload)
		return nil, errors.New(" lead creation failed")
	}

	res := api.CreateLeadResponse{
		NbfcLeadId:      fmt.Sprintf("%s", response.Results.ApplicationId),
		Status:          "SUCCESS",
		NbfcBorrowerId:  fmt.Sprintf("%d", response.Results.CustomerId),
		RejectionReason: map[string]string{},
		RedirectUrl:     response.Results.RedirectUrl,
	}
	return &res, nil
}

func (lp ApiClient) toCreateLeadPayload(req *api.CreateLeadRequest) (*initiateApplicationPayload, error) {
	fName, mName, lName := utils.GetFirstMiddleLastName(req.GetUserPersonalDetails().Name)

	pincode, err := strconv.Atoi(req.UserPersonalDetails.Pincode)
	if err != nil {
		return nil, err
	}

	return &initiateApplicationPayload{
		FirstName:      fName,
		MiddleName:     mName,
		LastName:       lName,
		Gender:         req.GetUserPersonalDetails().Gender,
		Dob:            req.GetUserPersonalDetails().Dob,
		PanNumber:      req.GetUserPersonalDetails().Pan,
		Pincode:        pincode,
		Income:         int(req.GetUserBureauDetails().EstimatedIncome),
		EmploymentType: dummyOccupation,
		Phone:          req.GetUserPersonalDetails().PhoneNumber,
		Email:          fmt.Sprintf("%s@rmnkmr.com", req.GetUserPersonalDetails().PhoneNumber),
	}, nil
}

func checkDuplicate(ctx context.Context, lp ApiClient, request *api.CreateLeadRequest) (bool, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, check_uplicate)
	method := "POST"

	payload := &checkDuplicatePayload{
		Email: fmt.Sprintf("%s@rmnkmr.com", request.UserPersonalDetails.PhoneNumber),
		Phone: request.UserPersonalDetails.PhoneNumber,
	}
	apiResponse := &checkDuplicateResponse{}
	err := httpCall(ctx, lp, url, method, payload, apiResponse)
	if err != nil {
		log.Error(ctx, err, "checkDuplicate: API call failed", "err", err, "req", payload, "apiResponse", apiResponse)
		return false, err
	}

	if apiResponse == nil || !apiResponse.Status {
		log.Error(ctx, err, "checkDuplicate: API call failed", "err", err, "req", payload, "apiResponse", apiResponse)
		return false, errors.New("checkDuplicate: API call failed")
	}
	return apiResponse.Results, nil
}
