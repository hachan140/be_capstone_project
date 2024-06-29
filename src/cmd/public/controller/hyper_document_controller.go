package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"github.com/gin-gonic/gin"
)

type HyperDocumentController struct {
	hyperDocumentService services.IHyperDocumentService
}

func NewHyperDocumentController(hyperDocumentService services.IHyperDocumentService) HyperDocumentController {
	return HyperDocumentController{hyperDocumentService: hyperDocumentService}
}

func (h *HyperDocumentController) FilterHyperDocument(ctx *gin.Context) {
	tag := "[HyperDocumentController] "
	documentFilterParams := request.HyperDocumentFilterParam{}
	if err := ctx.ShouldBindQuery(&documentFilterParams); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	res, err := h.hyperDocumentService.FilterHyperDocument(ctx, documentFilterParams)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandleCustomMessage(ctx, err.ServiceCode, err.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
}
