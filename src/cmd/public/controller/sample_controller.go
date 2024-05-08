package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SampleController struct {
	sampleService services.ISampleService
}

func NewSampleController(sampleService services.ISampleService) SampleController {
	return SampleController{sampleService: sampleService}
}

func (c *SampleController) CreateSampleController(ctx *gin.Context) {
	tag := "[CreateSampleController] "
	var req request.CreateSampleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to parse request body: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, common.ErrCodeInvalidRequest,
			fmt.Sprintf("Failed to bind the request's body to create sample request: %v", err))
		return
	}
	if err := req.Validate(); err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to parse request body: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, common.ErrCodeInvalidRequest,
			fmt.Sprintf(":Failed to bind the request's body to create sample request %v", err))
		return
	}
	if err := c.sampleService.CreateSampleService(ctx, &req); err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInternalError)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}
