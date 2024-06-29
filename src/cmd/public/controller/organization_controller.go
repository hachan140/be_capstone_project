package controller

import (
	"be-capstone-project/src/cmd/public/apihelper"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type OrganizationController struct {
	organizationService services.IOrganizationService
}

func NewOrganizationController(orgService services.IOrganizationService) OrganizationController {
	return OrganizationController{organizationService: orgService}
}

func (o *OrganizationController) CreateOrganization(ctx *gin.Context) {
	tag := "[CreateOrganizationController] "
	var req request.CreateOrganizationRequest
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
	err := o.organizationService.CreateOrganization(uint(userID), &req)
	if err != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (o *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	tag := "[UpdateOrganizationController] "
	orgIDRaw := ctx.Param("id")
	orgID, err := strconv.ParseUint(orgIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	var req request.UpdateOrganizationRequest
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
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	email := ""
	userEmail, ok := ctx.Get("email")
	if ok {
		email = fmt.Sprintf("%v", userEmail)
	}
	req.UpdatedBy = email
	errR := o.organizationService.UpdateOrganization(uint(orgID), uint(userID), &req)
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", errR)
		apihelper.AbortErrorHandle(ctx, errR.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
	return
}

func (o *OrganizationController) ViewOrganization(ctx *gin.Context) {
	tag := "[ViewOrganizationController] "
	orgIDRaw := ctx.Param("id")
	orgID, err := strconv.ParseUint(orgIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	res, errR := o.organizationService.FindOrganizationByID(uint(orgID), uint(userID))
	if errR != nil {
		logger.ErrorCtx(ctx, tag+"Failed to create sample with error: %v", errR)
		apihelper.AbortErrorHandle(ctx, errR.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, res)
	return
}

func (o *OrganizationController) AddPeopleToOrganization(ctx *gin.Context) {
	tag := "[AddPeopleToOrganizationController] "
	orgIDRaw := ctx.Param("id")
	orgID, err := strconv.ParseUint(orgIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userIDRaw, _ := ctx.Get("user_id")
	userID, _ := strconv.ParseUint(userIDRaw.(string), 10, 32)
	var req request.AddPeopleToOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	validEmails, errRes := o.organizationService.AddPeopleToOrganization(ctx, uint(orgID), uint(userID), req.Emails)
	if errRes != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, errRes.ServiceCode)
		return
	}
	res := map[string][]string{
		"valid_email": validEmails,
	}
	apihelper.SuccessfulHandle(ctx, res)
}

func (o *OrganizationController) AcceptOrganizationInvitation(ctx *gin.Context) {
	tag := "[AcceptOrganizationInvitationController] "
	orgIDRaw := ctx.Param("orgID")
	orgID, err := strconv.ParseUint(orgIDRaw, 10, 32)
	if err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, common.ErrCodeInvalidRequest)
		return
	}
	userEmail := ctx.Param("userEmail")
	if err := o.organizationService.AcceptOrganizationInvitation(uint(orgID), userEmail); err != nil {
		logger.Error(ctx, tag, err)
		apihelper.AbortErrorHandle(ctx, err.ServiceCode)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
}
