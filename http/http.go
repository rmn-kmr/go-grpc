package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/rmnkmr/lsp/log"

	errors "github.com/rmnkmr/go-common/errors"
)

// http request
type Request struct {
	Url         string
	BearerToken string
	Body        interface{}
	Headers     map[string]string
}

// http response
type Response struct {
	StatusCode int
	Body       []byte
}

// http post request
func (req *Request) Post(ctx context.Context) (*Response, error) {
	jsonBody, err := json.Marshal(req.Body)
	if err != nil {
		log.Error(ctx, err, "Http Request Lib: Error marshalling request body", "Req", req)
		return nil, err
	}
	httpReq, err := http.NewRequest(http.MethodPost, req.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(ctx, err, "Http Request Lib: http post request creation failed:", "Req", req)
		return nil, errors.Internal()
	}
	return req.ProcessRequest(ctx, httpReq)
}

// http get request
func (req *Request) Get(ctx context.Context) (*Response, error) {
	httpReq, err := http.NewRequest(http.MethodGet, req.Url, nil)
	if err != nil {
		log.Error(ctx, err, "Http Request Lib: http get request creation failed:", "Req", req)
		return nil, errors.Internal()
	}
	return req.ProcessRequest(ctx, httpReq)
}

// http request response process
func (resp *Request) ProcessRequest(ctx context.Context, httpReq *http.Request) (*Response, error) {
	httpReq.Header.Add("Content-Type", "application/json")
	// add bearer token if present
	if resp.BearerToken != "" {
		httpReq.Header.Add("Authorization", "Bearer "+resp.BearerToken)
	}

	// To Add Non-Conventional Headers
	for key, value := range resp.Headers {
		httpReq.Header.Add(key, value)
	}

	httpResp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		log.Error(ctx, err, "Http Request Lib: http request do failed:", "Req", httpReq)
		return nil, errors.Internal()
	}
	defer httpResp.Body.Close()
	body, _ := ioutil.ReadAll(httpResp.Body)

	log.Debug(ctx, "Http Request Lib:  raw response body:", "method", httpReq.Method, "Resp Body", string(body))
	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
	}, nil
}
