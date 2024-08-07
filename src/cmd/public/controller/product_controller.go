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

type ProductController struct {
	productService services.IProductService
}

func NewProductController(productService services.IProductService) ProductController {
	return ProductController{productService: productService}
}

func (p *ProductController) CreateProduct(ctx *gin.Context) {
	tag := "[CreateProduct] "
	var req request.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	err := p.productService.CreateProduct(ctx, uint(userID), &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create product with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, err.ServiceCode, err.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (p *ProductController) GetProductByID(ctx *gin.Context) {
	tag := "[GetProductByID] "
	prodIDRaw := ctx.Param("id")
	prodID, err := strconv.ParseUint(prodIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	res, errR := p.productService.GetProductByID(ctx, uint(prodID))
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to get product with error: %v", errR)
		apihelper.AbortErrorHandleCustomMessage(ctx, errR.ServiceCode, errR.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}

func (p *ProductController) SearchProductByName(ctx *gin.Context) {
	tag := "[ViewCategoryController] "
	name := ctx.Param("name")
	res, errG := p.productService.GetProductsByName(ctx, name)
	if errG != nil {
		logger.ErrorCtx(ctx, tag+"Failed to get category with error: %v", errG.Message)
		apihelper.AbortErrorHandleCustomMessage(ctx, errG.ServiceCode, errG.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}

func (p *ProductController) UpdateProduct(ctx *gin.Context) {
	tag := "[UpdateProduct] "
	prodIDRaw := ctx.Param("id")
	prodID, err := strconv.ParseUint(prodIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	var req request.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	errR := p.productService.UpdateProduct(ctx, uint(userID), uint(prodID), &req)
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to update product with error: %v", errR)
		apihelper.AbortErrorHandleCustomMessage(ctx, errR.ServiceCode, errR.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (p *ProductController) ViewListProduct(ctx *gin.Context) {
	tag := "[ViewListProduct] "
	productParam := request.ListProductRequest{}
	if err := ctx.ShouldBindQuery(&productParam); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	res, errR := p.productService.GetListProducts(ctx, productParam.Page, productParam.PageSize)
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to get list categories with error: %v", errR)
		apihelper.AbortErrorHandleCustomMessage(ctx, errR.ServiceCode, errR.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}

func (p *ProductController) DeleteProduct(ctx *gin.Context) {
	tag := "[DeleteProduct] "
	prodIDRaw := ctx.Param("id")
	prodID, err := strconv.ParseUint(prodIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	errR := p.productService.DeleteProduct(ctx, uint(userID), uint(prodID))
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to delete product with error: %v", errR)
		apihelper.AbortErrorHandleCustomMessage(ctx, errR.ServiceCode, errR.Message)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}
