package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CategoryController struct {
	categoryService services.ICategoryService
}

func NewCategoryController(categoryService services.ICategoryService) CategoryController {
	return CategoryController{categoryService: categoryService}
}

func (o *CategoryController) CreateCategory(ctx *gin.Context) {
	tag := "[CreateCategoryController] "
	var req request.CreateCategoryRequest
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
	email := ""
	userEmail, ok := ctx.Get("email")
	if ok {
		email = fmt.Sprintf("%v", userEmail)
	}
	req.CreatedBy = email
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	err := o.categoryService.CreateCategory(ctx, uint(userID), &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create category with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (o *CategoryController) UpdateCategory(ctx *gin.Context) {
	tag := "[UpdateCategoryController] "
	catIDRaw := ctx.Param("id")
	catID, err := strconv.ParseUint(catIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	var req request.UpdateCategoryRequest
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
	email := ""
	userEmail, ok := ctx.Get("email")
	if ok {
		email = fmt.Sprintf("%v", userEmail)
	}
	req.UpdatedBy = email
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	errRes := o.categoryService.UpdateCategoryByID(ctx, uint(userID), uint(catID), &req)
	if errRes != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create category with error: %v", errRes)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, errRes.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (c *CategoryController) ViewCategoryByID(ctx *gin.Context) {
	tag := "[ViewCategoryController] "
	catRaw := ctx.Param("id")
	catID, err := strconv.ParseUint(catRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}

	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)

	res, err := c.categoryService.GetCategoryByID(ctx, uint(catID), uint(userID))
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to get category with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}

func (c *CategoryController) ViewListCategoryByOrganization(ctx *gin.Context) {
	tag := "[ViewCategoryController] "
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	orgIDRaw := ctx.Param("id")
	orgID, err := strconv.ParseUint(orgIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	var req request.GetListCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	res, err := c.categoryService.ListCategories(ctx, uint(orgID), uint(userID), &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to get list categories with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}
