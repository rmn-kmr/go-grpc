package stashfin

// Routes

const (
	loginAPI             = "/v3/login-client"
	initiate_application = "/v3/initiate-application"
	check_uplicate       = "/v3/check-duplicate"
	check_status         = "/v3/check-status?application_id="
	check_limit_status   = "/v3/check-limit-status?application_id="
)
const (
	Permanent_Address         = 1
	Current_Residence_Address = 2
	Marital_Status_single     = 1
	Marital_Status_Married    = 2
	dummyOccupation           = 1 // salaried
)

const (
	AADHAR_KEY    = "aadhar"
	Self_Employed = "Self Employed"
)

const (
	PASSED   = "Passed"
	REJECTED = "Rejected"
	HOLD     = "hold"
)

type updateDocumentID int32

const (
	updateDocumentID_aadhaar updateDocumentID = 1
)

type updateDocumentType string

const (
	updateDocumentTyoe_aadhaar updateDocumentType = "aadhaar"
)

const AvgNumberOfDaysInMonth = 30.44
