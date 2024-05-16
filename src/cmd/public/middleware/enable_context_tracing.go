package middleware

import (
	"be-capstone-project/src/cmd/public/config"
	"be-capstone-project/src/internal/core/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EnableCoreMiddlewareRequestTracing(e *gin.Engine, cfg config.Config) {
	var handlers []func(next http.Handler) http.Handler
	// Wrapping all the default golang lib net/http middleware into gin
	handlers = append(handlers,
		middleware.AdvancedResponseWriter(),
		middleware.RequestContext(cfg),
		middleware.CorrelationId(),
	)

	e.Use(WrapAll(handlers)...)
}
