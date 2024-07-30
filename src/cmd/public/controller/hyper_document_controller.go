package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

type HyperDocumentController struct {
	hyperDocumentService services.IHyperDocumentService
	searchService        services.ISearchService
}

func NewHyperDocumentController(hyperDocumentService services.IHyperDocumentService, searchService services.ISearchService) HyperDocumentController {
	return HyperDocumentController{hyperDocumentService: hyperDocumentService, searchService: searchService}
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

func (h *HyperDocumentController) SearchDocumentAndOrNot(ctx *gin.Context) {
	tag := "[HyperDocumentController] "
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	var req request.SearchAndOrNotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	documents, err := h.searchService.SearchDocumentAndOrNot(&req, uint(userID))
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandleCustomMessage(ctx, err.ServiceCode, err.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, documents)
}

func (h *HyperDocumentController) GetSearchHistoryKeywords(ctx *gin.Context) {
	tag := "[HyperDocumentController] "
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	var req request.SearchHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	res, err := h.searchService.GetSearchKeywords(uint(userID), &req)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandleCustomMessage(ctx, err.ServiceCode, err.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
}

func (h *HyperDocumentController) SaveSearchHistory(ctx *gin.Context) {
	tag := "[HyperDocumentController] "
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	var req request.SaveSearchHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	if err := req.Validate(); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandleCustomMessage(ctx, common.ErrCodeInvalidRequest, err.Error())
		return
	}
	req.UserID = uint(userID)
	err := h.searchService.SaveSearchHistory(&req)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandleCustomMessage(ctx, err.ServiceCode, err.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
}
