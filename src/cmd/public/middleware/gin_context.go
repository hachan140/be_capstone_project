package middleware

import (
	"be-capstone-project/src/internal/core/web/constant"
	"be-capstone-project/src/internal/core/web/context"
	"github.com/gin-gonic/gin"
)

// InitContext assign request_attribute to gin context
func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAttributes := context.GetOrCreateRequestAttributes(c.Request)
		requestAttributes.Mapping = c.FullPath()
		c.Set(constant.ContextReqAttribute, requestAttributes)
		c.Set(constant.ContextAcceptLanguage, requestAttributes.AcceptLanguage)
		c.Next()
	}
}
