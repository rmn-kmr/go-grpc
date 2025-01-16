package callback

//func (c *Callback) OnHTTPCallback(ctx context.Context, req *HTTPCallback) (*emptypb.Empty, error) {
//	log.Debug(ctx, "OnHTTPCallback, callback request received", "callbackRequest", req)
//	// filter the query params for the key source
//	var provider string
//	for _, queryParam := range req.QueryParams {
//		if queryParam.Key == "source" {
//			provider = queryParam.Values[0]
//		}
//	}
//	paymentProvider, err := c.App.Queries.GetLoanProviderBySlug(ctx, provider)
//	if err != nil {
//		log.Error(ctx, err, "OnHTTPCallback: GetLoanProviderBySlug")
//		return nil, err
//	}
//
//	client := lsp.ApiClient{
//		Provider:                      paymentProvider.Provider,
//		ProviderAccountId:             paymentProvider.ProviderAccountID,
//		ProviderID:                    paymentProvider.ProviderID,
//		ProviderBaseUrl:               paymentProvider.ApiBaseUrl,
//		ClientId:                      paymentProvider.ApiKey,
//		ClientSecret:                  paymentProvider.ApiSecret,
//		Environment:                   c.App.Environment,
//		PreExpirationTokenRefreshMins: paymentProvider.PreExpirationTokenRefreshMins,
//	}
//
//	_, callBackErr := client.OnCallBack(ctx, &paymentProviderApi.OnCallBackRequest{
//		RequestBody: req.Body,
//	})
//
//	if callBackErr != nil {
//		log.Error(ctx, callBackErr, "OnHTTPCallback: OnCallBack")
//		return nil, callBackErr
//	}
//	return nil, err
//}
