package apihelper

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/web/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AbortErrorHandle handle abort error
func AbortErrorHandle(c *gin.Context, code int) {
	language, existLang := c.Get(constant.ContextAcceptLanguage)
	if !existLang {
		language = common.Vietnamese
	}
	errorResponse := common.GetErrorResponse(code, language.(string))
	c.JSON(errorResponse.HTTPCode, &Response{
		Meta: Meta{
			Code:    errorResponse.ServiceCode,
			Message: errorResponse.Message,
		},
		Data: nil,
	})
}

// AbortErrorHandleCustomMessage handle abort with custom message
func AbortErrorHandleCustomMessage(c *gin.Context, code int, message string) {
	language, existLang := c.Get(constant.ContextAcceptLanguage)
	if !existLang {
		language = common.Vietnamese
	}
	errorResponse := common.GetErrorResponse(code, language.(string))
	c.JSON(errorResponse.HTTPCode, &Response{
		Meta: Meta{
			Code:    errorResponse.ServiceCode,
			Message: message,
		},
		Data: nil,
	})
}

// AbortErrorResponseHandle handle abort with error response
func AbortErrorResponseHandle(c *gin.Context, errorResponse *common.ErrorResponse) {
	var out common.ErrorResponse
	language, existLang := c.Get(constant.ContextAcceptLanguage)
	if !existLang {
		language = common.Vietnamese
	}
	resByLanguage := common.GetErrorResponse(errorResponse.ServiceCode, language.(string))
	out.HTTPCode = resByLanguage.HTTPCode
	if errorResponse.CustomMsgParam != nil {
		out.Message = fmt.Sprintf(resByLanguage.Message, *errorResponse.CustomMsgParam...)
	} else {
		out.Message = resByLanguage.Message
	}
	out.ServiceCode = resByLanguage.ServiceCode
	c.JSON(out.HTTPCode, &Response{
		Meta: Meta{
			Code:    out.ServiceCode,
			Message: out.Message,
		},
		Data: nil,
	})
}

// SuccessfulHandle handle successful response
func SuccessfulHandle(c *gin.Context, data interface{}) {

	c.JSON(http.StatusOK, &Response{
		Meta: Meta{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
		Data: data,
	})
}
