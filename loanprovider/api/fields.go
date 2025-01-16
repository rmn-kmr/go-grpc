package api

import "time"

type AuthToken struct {
	ProviderID     string    `json:"provider_id"`
	AuthToken      string    `json:"auth_token"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
}

type InitiateLoanRequest struct {
	Amount int64 `json:"amount"`
}

type InitiateLoanResponse struct {
	LoanId string `json:"loan_id"`
}

type OnCallBackRequest struct {
	RequestBody string `json:"request_body"`
}

type CallbackResponse struct {
	Event string `json:"event"`
}
