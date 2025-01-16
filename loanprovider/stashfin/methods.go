package stashfin

import (
	"context"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/rmnkmr/lsp/log"
	"github.com/rmnkmr/lsp/utils"
	"math"
)

func ConvertAadhaarFileForStashfin(ctx context.Context, file []byte, fileName string, password string) ([]byte, error) {
	// unzip file
	unzippedXmlData, err := unZipFile(file, password)
	if err != nil {
		log.Error(ctx, err, "ConvertAadhaarFileForStashfin: unzipFile failed")
		return nil, err
	}
	// extract xml data and store in OkAadhaarKycData struct
	aadhaarKycData, err := parseAadhaarXML(unzippedXmlData)
	if err != nil {
		log.Error(ctx, err, "ConvertAadhaarFileForStashfin: parseAadhaarXML failed")
		return nil, err
	}
	// create XML data of type SfAadhaarKycData
	xmlData, err := createAadhaarXMLFileForStashfin(aadhaarKycData)
	if err != nil {
		log.Error(ctx, err, "ConvertAadhaarFileForStashfin: createAadhaarXMLFileForStashfin failed")
		return nil, err
	}
	return xmlData, nil
}

func httpCall(ctx context.Context, lp ApiClient, url string, method string, payload interface{}, response interface{}) error {
	authToken, err := lp.GetTokenFromProvider(ctx)
	if err != nil {
		log.Error(ctx, err, "httpCall: GetTokenFromProvider failed", "err", err, "authToken", authToken)
		return err
	}

	headers := map[string]string{
		"client-token": authToken.AuthToken,
		"Content-Type": utils.IfThenElse(method == http.MethodGet, "", utils.HttpRequestContentTypeJson).(string),
	}

	httpReq := &utils.HttpRequest{
		Endpoint:     url,
		Method:       method,
		Headers:      headers,
		Req:          payload,
		Resp:         response,
		Retry:        true,
		ResponseType: utils.HttpRequestContentTypeJson,
	}

	resp, err := utils.MakeHttpCall(
		ctx,
		httpReq)

	err = mapstructure.Decode(resp, response)
	if err != nil {
		log.Error(ctx, err, "httpCall: mapstructure.Decode failed", "err", err, "resp", resp)
		return err
	}

	return nil
}

func CalculateMonthlyInterest(interestPercentage float32) float32 {
	monthlyInterest := interestPercentage / 12
	return float32(int(monthlyInterest*100+0.5)) / 100.0
}

func CalculateMonthsForNumberOfDays(days int64) int {
	daysInMonth := AvgNumberOfDaysInMonth
	months := float64(days) / daysInMonth
	return int(math.Ceil(months))
}
