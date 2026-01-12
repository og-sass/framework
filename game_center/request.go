package game_center

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/og-saas/framework/utils/httpc"
)

func postRequest(ctx context.Context, config CenterConfig, url, currency string, body any) (*resty.Response, error) {
	resp, err := httpc.Do(ctx).
		SetBasicAuth(config.GetCurrencyConf(currency).Username, config.GetCurrencyConf(currency).Password).
		SetBody(body).
		Post(config.RequestURL + url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("POST request failed with status " + resp.Status() + ": " + resp.String())
	}
	return resp, nil
}

func getRequest(ctx context.Context, config CenterConfig, url, currency string, params url.Values) (*resty.Response, error) {
	resp, err := httpc.Do(ctx).
		SetBasicAuth(config.GetCurrencyConf(currency).Username, config.GetCurrencyConf(currency).Password).
		SetQueryParamsFromValues(params).
		Get(config.RequestURL + url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK { // 新增：检查HTTP状态码
		return nil, errors.New("GET request failed with status " + resp.Status() + ": " + resp.String())
	}
	return resp, nil
}
