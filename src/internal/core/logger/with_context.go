package logger

import (
	"be-capstone-project/src/internal/core/web/constant"
	webcontext "be-capstone-project/src/internal/core/web/context"
	"context"
)

type LoggingContext struct {
	CorrelationId string `json:"request_id,omitempty"`
	//UserId            string `json:"jwt_subject,omitempty"`
	DeviceId        string `json:"device_id,omitempty"`
	DeviceSessionId string `json:"device_session_id,omitempty"`
	//TechnicalUsername string `json:"technical_username,omitempty"`
}

func BuildLoggingContextFromReqAttr(requestAttributes *webcontext.RequestAttributes) *LoggingContext {
	return &LoggingContext{
		DeviceId:        requestAttributes.DeviceId,
		DeviceSessionId: requestAttributes.DeviceSessionId,
		CorrelationId:   requestAttributes.CorrelationId,
		//UserId:            requestAttributes.SecurityAttributes.UserId,
		//TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}
}

func keysAndValuesFromContext(ctx context.Context) []interface{} {
	if requestAttributes := webcontext.GetRequestAttributes(ctx); requestAttributes != nil {
		return []interface{}{constant.ContextReqMeta, BuildLoggingContextFromReqAttr(requestAttributes)}
	}
	return nil
}

// InfoCtx uses fmt.Sprint to construct and log a message.
func InfoCtx(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Infow(keysAndValuesFromContext(ctx), msgFormat, args...)
}

// DebugCtx uses fmt.Sprint to construct and log a message.
func DebugCtx(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Debugw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

// WarnCtx uses fmt.Sprint to construct and log a message.
func WarnCtx(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Warnw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

// ErrorCtx uses fmt.Sprint to construct and log a message. ErrorCtx level log with a stack trace
func ErrorCtx(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Errorw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

// FatalCtx uses fmt.Sprint to construct and log a message, then calls os.Exit.
func FatalCtx(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Fatalw(keysAndValuesFromContext(ctx), msgFormat, args...)
}
