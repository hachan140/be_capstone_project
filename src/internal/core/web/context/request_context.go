package context

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/web/constant"
	"context"
	"net/http"
)

func GetRequestAttributes(ctx context.Context) *RequestAttributes {
	reqAttrCtxValue := ctx.Value(constant.ContextReqAttribute)
	if reqAttrCtxValue == nil {
		return nil
	}
	requestAttributes, ok := reqAttrCtxValue.(*RequestAttributes)
	if !ok {
		return nil
	}
	return requestAttributes
}

func GetRequestAcceptLanguage(ctx context.Context) string {
	reqCtxAcceptLanguage := ctx.Value(constant.ContextAcceptLanguage)
	if reqCtxAcceptLanguage != nil {
		return common.Vietnamese
	}

	reqAcceptLanguage, ok := reqCtxAcceptLanguage.(string)
	if !ok {
		return common.Vietnamese
	}
	return reqAcceptLanguage
}

func GetOrCreateRequestAttributes(r *http.Request) *RequestAttributes {
	requestAttributes := GetRequestAttributes(r.Context())
	if requestAttributes == nil {
		return replaceNewRequestAttributes(r)
	}
	return requestAttributes
}

func replaceNewRequestAttributes(r *http.Request) *RequestAttributes {
	requestAttributes := makeRequestAttributes(r)
	*r = *r.WithContext(context.WithValue(r.Context(), constant.ContextReqAttribute, requestAttributes))
	return requestAttributes
}

func makeRequestAttributes(r *http.Request) *RequestAttributes {
	return &RequestAttributes{
		AcceptLanguage:  r.Header.Get(constant.HeaderAcceptLanguage),
		Uri:             r.URL.Path,
		Query:           r.URL.RawQuery,
		Url:             r.URL.String(),
		Method:          r.Method,
		CallerId:        getServiceClientName(r),
		DeviceId:        r.Header.Get(constant.HeaderDeviceId),
		DeviceSessionId: r.Header.Get(constant.HeaderDeviceSessionId),
		DeviceUuid:      r.Header.Get(constant.HeaderDeviceUuid),
		ClientIpAddress: getClientIpAddress(r),
		UserAgent:       r.Header.Get(constant.HeaderUserAgent),
		XForwardedFor:   r.Header.Get(constant.HeaderXForwardedFor),
		TrueClientIp:    r.Header.Get(constant.HeaderTrueClientIP),
		CfConnectingIP:  r.Header.Get(constant.HeaderCFConnectingIP),
		XChannel:        r.Header.Get(constant.HeaderXChannel),
	}
}

func getClientIpAddress(r *http.Request) string {
	if clientIpAddress := r.Header.Get(constant.HeaderClientIpAddress); len(clientIpAddress) > 0 {
		return clientIpAddress
	}
	return r.RemoteAddr
}

func getServiceClientName(r *http.Request) string {
	serviceName := r.Header.Get(constant.HeaderServiceClientName)
	if len(serviceName) > 0 {
		return serviceName
	}
	return ""
}
