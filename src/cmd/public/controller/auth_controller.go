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

func (a *AuthController) VerifyEmail(ctx *gin.Context) {
	tag := "[VerifyEmailController]"
	req := request.VerifyEmail{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}

	err := a.userService.UpdateUserStatusWhenEmailVerified(ctx, req.Email)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
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

func (a *AuthController) SocialLogin(ctx *gin.Context) {
	tag := "[LoginSocialController] "
	var req request.SocialLoginRequest
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
	res, err := a.userService.LoginSocial(ctx, &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to login social with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
}

func (a *AuthController) RefreshToken(ctx *gin.Context) {
	tag := "[LoginController] "
	var req request.RefreshTokenRequest
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
	res, err := a.userService.RefreshToken(ctx, &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to login social with error: %v", err)
		apihelper.AbortErrorHandleCustomMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
}

func (a *AuthController) ResetPasswordRequest(ctx *gin.Context) {
	tag := "[AuthController] "
	var req request.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	if err := req.Validate(); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	token, err := a.userService.ResetPasswordRequest(ctx, &req)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, token)
}

func (a *AuthController) ResetPassword(ctx *gin.Context) {
	tag := "[AuthController] "
	var req request.ResetPassword
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	if err := req.Validate(); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	email := ""
	userEmail, ok := ctx.Get("email")
	if ok {
		email = fmt.Sprintf("%v", userEmail)
	}
	if err := a.userService.ResetPassword(ctx, email, &req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
}
