package context

import "time"

type RequestAttributes struct {
	AcceptLanguage  string        `json:"accept_language"`
	ServiceCode     string        `json:"service_code"`
	StatusCode      int           `json:"status"`
	ExecutionTime   time.Duration `json:"duration_ms"`
	Uri             string        `json:"uri"`
	Query           string        `json:"query"`
	Mapping         string        `json:"mapping"`
	Url             string        `json:"url"`
	Method          string        `json:"method"`
	CallerId        string        `json:"caller_id"`
	DeviceId        string        `json:"device_id"`
	DeviceSessionId string        `json:"device_session_id"`
	DeviceUuid      string        `json:"device_uuid"`
	CorrelationId   string        `json:"correlation_id"`
	ClientIpAddress string        `json:"client_ip_address"`
	UserAgent       string        `json:"user_agent"`
	XForwardedFor   string        `json:"x_forwarded_for"`
	TrueClientIp    string        `json:"true_client_ip"`
	CfConnectingIP  string        `json:"cf_connecting_ip"`
	XChannel        string        `json:"x_channel"`
}
