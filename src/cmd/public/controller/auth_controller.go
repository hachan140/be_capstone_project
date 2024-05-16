package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	userService services.IUserService
}

func NewAuthController(userService services.IUserService) AuthController {
	return AuthController{userService: userService}
}

func (a *AuthController) Signup(ctx *gin.Context) {
	tag := "[SignupController]"
	var req request.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "", err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	if err := req.Validate(); err != nil {
		logger.Error(ctx, "", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, common.ErrCodeInvalidRequest, err.Error())
		return
	}
	err := a.userService.CreateUser(ctx, &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (a *AuthController) Login(ctx *gin.Context) {
	tag := "[LoginController] "
	var req request.LoginRequest
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
	res, err := a.userService.LoginByUserEmail(ctx, &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to login with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
}
