package session

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/web/context"
	"github.com/gin-gonic/gin"
)

func GetReqAcceptLanguage(c *gin.Context) common.ResponseLanguage {
	lang := context.GetRequestAcceptLanguage(c)
	return common.StringToLanguage(lang)
}
