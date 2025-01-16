package stashfin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rmnkmr/go-common/errors"
	api2 "github.com/rmnkmr/lsp/loanprovider/api"
	"github.com/rmnkmr/lsp/log"
)

type LoginPayload struct {
	Id           string `json:"id"`
	ClientSecret string `json:"client_secret"`
}

type LoginResponse struct {
	Status  bool   `json:"status"`
	Results string `json:"results"`
}

func (lp ApiClient) login(ctx context.Context) (*api2.AuthToken, error) {
	url := fmt.Sprintf("%s%s", lp.ProviderBaseUrl, loginAPI)
	method := "POST"

	payload := &LoginPayload{
		Id:           lp.ClientId,
		ClientSecret: lp.ClientSecret,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))

	if err != nil {
		log.Error(ctx, err, "login: login failed", "err", err, "req", req)
		return nil, nil
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(ctx, err, "login: API call failed", "err", err, "req", req)
		return nil, nil
	}
	defer res.Body.Close()

	llResp := LoginResponse{}

	err = json.NewDecoder(res.Body).Decode(&llResp)
	if err != nil {
		log.Error(ctx, err, "login: failed to parse response body: %v", err)
		return nil, errors.Internal()
	}

	if !llResp.Status {
		log.Error(ctx, err, "login: login failed", "llResp", llResp, "req", req)
		return nil, errors.New(" login failed")
	}

	return &api2.AuthToken{
		AuthToken:  llResp.Results,
		CreateTime: time.Now(),
	}, nil
}
