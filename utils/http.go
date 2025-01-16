package utils

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rmnkmr/go-common/encoding/json"
	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/go-common/hack/auth"
	"github.com/rmnkmr/lsp/log"
)

type HttpRequest struct {
	Endpoint     string
	Method       string
	Headers      map[string]string
	Req          interface{}
	Resp         interface{}
	Retry        bool
	ResponseType string
}

func MakeHttpCall(ctx context.Context, httpReq *HttpRequest) (interface{}, error) {
	log.Info(ctx, "MakeHttpCall: received request", "Req", httpReq.Req, "endpoint", httpReq.Endpoint)

	var payload io.Reader
	var err error

	if val, found := httpReq.Headers[ContentType]; found && val == HttpRequestContentTypeJson {
		jsonBody, err := json.Marshal(httpReq.Req)
		if err != nil {
			log.Error(ctx, err, "MakeHttpCall: json marshaling failed", "Req", httpReq.Req)
			return nil, err
		}
		payload = bytes.NewBuffer(jsonBody)
	} else if val, found := httpReq.Headers[ContentType]; found && val == HttpFormURlEncodedType {
		payload = strings.NewReader(httpReq.Req.(string))
	} else if val, found := httpReq.Headers[ContentType]; found && val == HttpRequestContentTypeXML {
		var ok bool
		byteArray, ok := httpReq.Req.([]byte)
		if !ok {
			log.Error(ctx, errors.Internal(), "MakeHttpCall: invalid type", "Req", httpReq.Req)
			return nil, errors.Internal()
		}

		payload = bytes.NewReader(byteArray)
	}

	var httpResp *http.Response
	if httpReq.Retry {
		httpResp, err = retryableHttpCall(ctx, httpReq.Req, httpReq.Method, httpReq.Endpoint, httpReq.Headers, payload)
		if err != nil {
			log.Error(ctx, err, "MakeHttpCall: retryableHttpCall failed", "Req", httpReq.Req)
			return nil, err
		}
	} else {
		httpResp, err = nonRetryableHttpCall(ctx, httpReq.Req, httpReq.Method, httpReq.Endpoint, httpReq.Headers, payload)
		if err != nil {
			log.Error(ctx, err, "MakeHttpCall: nonRetryableHttpCall failed", "Req", httpReq.Req)
			return nil, err
		}
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode > 300 {
		// Check for error codes
		log.Info(ctx, "MakeHttpCall: request failed",
			"statusCode", httpResp.StatusCode, "status", httpResp.Status)
		if httpReq.Resp != nil {
			if httpReq.ResponseType == HttpRequestContentTypeXML {
				body, err := ioutil.ReadAll(httpResp.Body)
				if err != nil {
					log.Error(ctx, err, "MakeHttpCall: failed to read failure response body",
						"Req", httpReq.Req)
					return nil, err
				}
				log.Error(ctx, err, "MakeHttpCall: response", "failureResp", string(body))
				return body, err
			}
			err = json.NewDecoder(httpResp.Body).Decode(&httpReq.Resp)
			if err != nil {
				log.Error(ctx, err, "MakeHttpCall: failed to parse failure response body",
					"Req", httpReq.Req)
				return httpReq.Resp, err
			}
			log.Error(ctx, err, "MakeHttpCall: response", "Resp", httpReq.Resp)
			return httpReq.Resp, nil
		}
		log.Info(ctx, "MakeHttpCall: response", "Resp", httpReq.Resp)
		return httpReq.Resp, nil
	}

	// 3. Read http response
	// If response type is XML
	if httpReq.ResponseType == HttpRequestContentTypeXML {
		if err := xml.NewDecoder(httpResp.Body).Decode(httpReq.Resp); err != nil {
			log.Error(ctx, err, "MakeHttpCall: failed to parse failure response body",
				"Req", httpReq.Req)
			return nil, err
		}
		log.Info(ctx, "MakeHttpCall: response", "Resp", httpReq.Resp)
		return httpReq.Resp, err
	}

	// If response type is JSON
	err = json.NewDecoder(httpResp.Body).Decode(&httpReq.Resp)
	if err != nil {
		log.Error(ctx, err, "MakeHttpCall: failed to parse response body",
			"Req", httpReq.Req)
		return nil, err
	}
	log.Info(ctx, "MakeHttpCall: response", "Resp", httpReq.Resp)
	return httpReq.Resp, nil
}

func retryableHttpCall(ctx context.Context,
	req interface{},
	method, endpoint string,
	headers map[string]string,
	payload io.Reader) (*http.Response, error) {

	// Adding a retryable http client
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = HttpRetryAttempts
	httpClient.RetryWaitMin = HttpMinRetrySeconds * time.Second
	httpClient.RetryWaitMax = HttpMaxRetrySeconds * time.Second
	httpClient.CheckRetry = retryablehttp.ErrorPropagatedRetryPolicy

	log.Info(ctx, "retryableHttpCall: request", "Req", req)

	// 1. Prepare http request
	httpReq, err := retryablehttp.NewRequest(method, endpoint, payload)
	if err != nil {
		log.Error(ctx, err, "retryableHttpCall: Error while creating HTTP request", "Req", req)
		return nil, err
	}

	for k, v := range headers {
		httpReq.Header.Add(k, v)
	}

	// 2. Hit external api
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, err, "retryableHttpCall: Error while hitting API", "Req", req)
		return nil, err
	}

	return httpResp, nil
}

func nonRetryableHttpCall(ctx context.Context,
	req interface{},
	method, endpoint string,
	headers map[string]string,
	payload io.Reader) (*http.Response, error) {

	httpClient := auth.NewHttpClientV2()

	log.Info(ctx, "nonRetryableHttpCall: request", "Req", req)

	// 1. Prepare http request
	httpReq, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		log.Error(ctx, err, "nonRetryableHttpCall: Error while creating HTTP request", "Req", req)
		return nil, err
	}

	for k, v := range headers {
		httpReq.Header.Add(k, v)
	}

	// 2. Hit external api
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, err, "nonRetryableHttpCall: Error while hitting API", "Req", req)
		return nil, err
	}

	return httpResp, nil

}
