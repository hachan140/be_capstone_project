package client

import (
	"be-capstone-project/src/internal/core/common_configs"
	"be-capstone-project/src/internal/core/web/constant"
	webContext "be-capstone-project/src/internal/core/web/context"
	"context"
	"net/http"
)

type TraceableHttpClient struct {
	appConfig common_configs.App
	client    HttpClient
}

func NewTraceableHttpClient(appConfig common_configs.App) ContextualHttpClient {
	c := http.Client{}
	defaultClient := NewDefaultHttpClient(&c)

	return &TraceableHttpClient{
		appConfig: appConfig,
		client:    defaultClient,
	}
}

func (t *TraceableHttpClient) Get(ctx context.Context, url string, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodGet, url, nil, result, options...)
}

func (t *TraceableHttpClient) Post(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPost, url, body, result, options...)
}

func (t *TraceableHttpClient) Put(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPut, url, body, result, options...)
}

func (t *TraceableHttpClient) Patch(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPatch, url, body, result, options...)
}

func (t *TraceableHttpClient) Delete(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodDelete, url, body, result, options...)
}

func (t *TraceableHttpClient) Request(ctx context.Context, method string, url string, body interface{},
	result interface{}, options ...RequestOption) (*HttpResponse, error) {
	httpOpts := []RequestOption{
		WithHeader(constant.HeaderServiceClientName, t.appConfig.Name),
	}
	if reqAttrs := webContext.GetRequestAttributes(ctx); reqAttrs != nil {
		httpOpts = append(httpOpts,
			WithHeader(constant.HeaderCorrelationId, reqAttrs.CorrelationId),
			WithHeader(constant.HeaderDeviceId, reqAttrs.DeviceId),
			WithHeader(constant.HeaderDeviceSessionId, reqAttrs.DeviceSessionId),
			WithHeader(constant.HeaderClientIpAddress, reqAttrs.ClientIpAddress),
			WithHeader(constant.HeaderUserAgent, reqAttrs.UserAgent),
		)
	}
	httpOpts = append(httpOpts, options...)
	return t.client.Request(method, url, body, result, httpOpts...)
}
