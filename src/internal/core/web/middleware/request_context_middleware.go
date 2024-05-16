package middleware

import (
	"be-capstone-project/src/cmd/public/config"
	"be-capstone-project/src/internal/core/logger"
	"be-capstone-project/src/internal/core/web/constant"
	"be-capstone-project/src/internal/core/web/context"
	"errors"
	"net/http"
	"time"
)

// RequestContext middleware responsible to inject attributes to the request's context.
// This middleware should be run as soon as possible to
// create a uniform context for the request.
func RequestContext(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestAttributes := context.GetOrCreateRequestAttributes(r)
			requestAttributes.ServiceCode = cfg.App.Name
			next.ServeHTTP(w, r)
			if advancedResponseWriter, err := getAdvancedResponseWriter(w); err != nil {
				logger.WarnCtx(r.Context(), "Cannot detect AdvancedResponseWriter with error [%s]", err.Error())
			} else {
				requestAttributes.StatusCode = advancedResponseWriter.Status()
			}
			requestAttributes.ExecutionTime = time.Now().Sub(start)
			logger.Infow(
				[]interface{}{constant.ContextReqAttribute, requestAttributes},
				"finish router")
		})
	}
}

func getAdvancedResponseWriter(w http.ResponseWriter) (*context.AdvancedResponseWriter, error) {
	if advancedResponseWriter, ok := w.(*context.AdvancedResponseWriter); ok {
		return advancedResponseWriter, nil
	}
	if wrappingWriter, ok := w.(context.WrappingResponseWriter); ok {
		if advancedResponseWriter, ok := wrappingWriter.Writer().(*context.AdvancedResponseWriter); ok {
			return advancedResponseWriter, nil
		}
		return nil, errors.New("ResponseWriter is wrapped by more than two level")
	}
	return nil, errors.New("your ResponseWriter is not implement context.WrappingResponseWriter")
}
