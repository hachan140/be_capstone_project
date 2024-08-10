package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/common_configs"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/logger"
	"be-capstone-project/src/internal/core/utils"
	"context"
	"errors"
	_ "gopkg.in/gomail.v2"
	"net/http"
	"time"
)

type IOrganizationService interface {
	CreateOrganization(userId uint, req *request.CreateOrganizationRequest) *common.ErrorCodeMessage
	UpdateOrganization(orgID uint, userID uint, req *request.UpdateOrganizationRequest) *common.ErrorCodeMessage
	UpdateOrganizationStatus(orgID uint, userID uint, req *request.UpdateOrganizationStatusRequest) *common.ErrorCodeMessage
	FindOrganizationByID(orgID uint, userID uint) (*dtos.Organization, *common.ErrorCodeMessage)
	CheckUserRoleInOrganization(orgID uint, userID uint) (bool, bool, error)
	AddPeopleToOrganization(ctx context.Context, orgID uint, userID uint, req *request.AddPeopleToOrganizationRequest) ([]string, *common.ErrorCodeMessage)
	RemoveUserFromOrganization(ctx context.Context, orgID uint, userID uint, req *request.RemoveUserFromOrganizationRequest) *common.ErrorCodeMessage
	AcceptOrganizationInvitation(orgID uint, deptID uint, userEmail string) *common.ErrorCodeMessage
	AssignPeopleTobeManager(ctx context.Context, orgID uint, userID uint, req *request.AssignPeopleToManagerRequest) *common.ErrorCodeMessage
	RecallPeopleTobeManager(ctx context.Context, orgID uint, userID uint, req *request.RecallPeopleManagerRequest) *common.ErrorCodeMessage
	CheckUserAlreadyRequestCreateOrganization(userId uint) (bool, *common.ErrorCodeMessage)
}

type OrganizationService struct {
	organizationRepository postgres.IOrganizationRepository
	userRepository         postgres.IUserRepository
	emailConfig            common_configs.EmailConfig
	categoryRepository     postgres.ICategoryRepository
	documentRepository     postgres.IDocumentRepository
}

func NewOrganizationService(orgRepo postgres.IOrganizationRepository, userRepository postgres.IUserRepository,
	emailConfig common_configs.EmailConfig, categoryRepository postgres.ICategoryRepository, documentRepository postgres.IDocumentRepository) IOrganizationService {
	return &OrganizationService{
		organizationRepository: orgRepo,
		userRepository:         userRepository,
		emailConfig:            emailConfig,
		categoryRepository:     categoryRepository,
		documentRepository:     documentRepository,
	}
}

func (o *OrganizationService) CreateOrganization(userId uint, req *request.CreateOrganizationRequest) *common.ErrorCodeMessage {
	user, err := o.userRepository.FinduserByID(userId)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID != 0 {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyInOtherOrganization,
			Message:     common.ErrMessageUserAlreadyInOtherOrganization,
		}
	}
	orgByName, err := o.organizationRepository.FindOrganizationByName(req.Name)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if orgByName != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationExisted,
			Message:     common.ErrMessageOrganizationExisted,
		}
	}
	orgByAuthor, err := o.organizationRepository.FindOrganizationByAuthor(user.Email)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if orgByAuthor != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyRequestCreateOrganization,
			Message:     common.ErrMessageUserAlreadyCreateOrganizationRequest,
		}
	}
	orgModel := &model.Organization{
		Name:        req.Name,
		Description: req.Description,
		IsOpenai:    req.IsOpenai,
		CreatedAt:   time.Now(),
		CreatedBy:   req.CreatedBy,
		Status:      0,
		LimitData:   common.LIMIT_DATA,
	}
	if err := o.organizationRepository.CreateOrganization(orgModel); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (o *OrganizationService) CheckUserAlreadyRequestCreateOrganization(userId uint) (bool, *common.ErrorCodeMessage) {
	user, err := o.userRepository.FinduserByID(userId)
	if err != nil {
		return false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID != 0 {
		return false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyInOtherOrganization,
			Message:     common.ErrMessageUserAlreadyInOtherOrganization,
		}
	}
	orgByAuthor, err := o.organizationRepository.FindOrganizationByAuthor(user.Email)
	if err != nil {
		return false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if orgByAuthor != nil && orgByAuthor.Status == 0 {
		return true, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyRequestCreateOrganization,
			Message:     common.ErrMessageUserAlreadyCreateOrganizationRequest,
		}
	}
	return false, nil
}

func (o *OrganizationService) UpdateOrganization(orgID uint, userID uint, req *request.UpdateOrganizationRequest) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	_, isOrganizationManager, err := o.CheckUserRoleInOrganization(orgID, userID)
	if !isOrganizationManager || err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if req.Name != nil && *req.Name != org.Name {
		orgByName, err := o.organizationRepository.FindOrganizationByName(*req.Name)
		if err != nil {
			return &common.ErrorCodeMessage{
				HTTPCode:    http.StatusInternalServerError,
				ServiceCode: common.ErrCodeInternalError,
				Message:     err.Error(),
			}
		}
		if orgByName != nil {
			return &common.ErrorCodeMessage{
				HTTPCode:    http.StatusBadRequest,
				ServiceCode: common.ErrCodeOrganizationExisted,
				Message:     common.ErrMessageOrganizationExisted,
			}
		}
	}

	orgToUpdate, err := o.buildOrganizationUpdateQuery(org, req)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := o.organizationRepository.UpdateOrganization(orgID, orgToUpdate); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (o *OrganizationService) UpdateOrganizationStatus(orgID uint, userID uint, req *request.UpdateOrganizationStatusRequest) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	_, isOrganizationManager, _ := o.CheckUserRoleInOrganization(orgID, userID)
	if !isOrganizationManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}

	oStatus, deptStatus, cStatus, docStatus := 0, 0, 0, 0

	if *req.Status == 1 {
		oStatus = 1
		deptStatus = 1
		cStatus = 1
		docStatus = 1
	}

	if *req.Status == 3 {
		oStatus = 3
		deptStatus = 2
		cStatus = 2
		docStatus = 2
	}

	if err := o.organizationRepository.UpdateOrganizationStatus(orgID, oStatus); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := o.categoryRepository.UpdateCategoriesStatusByOrganizationID(orgID, cStatus); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := o.categoryRepository.UpdateDepartmentsStatusByOrganizationID(orgID, deptStatus); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if err := o.documentRepository.UpdateDocumentStatusByOrganizationID(orgID, docStatus); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	return nil
}

func (o *OrganizationService) buildOrganizationUpdateQuery(orgExist *model.Organization, orgUpdate *request.UpdateOrganizationRequest) (*model.Organization, error) {
	var org model.Organization
	if orgUpdate == nil {
		return nil, errors.New(common.ErrMessageInvalidRequest)
	}
	if orgUpdate.Name != nil {
		org.Name = *orgUpdate.Name
	} else {
		org.Name = orgExist.Name
	}
	if orgUpdate.Description != nil {
		org.Description = *orgUpdate.Description
	} else {
		org.Description = orgExist.Description
	}
	if orgUpdate.Status != nil {
		org.Status = *orgUpdate.Status
	}
	org.UpdatedAt = time.Now()
	org.UpdatedBy = orgUpdate.UpdatedBy
	org.CreatedBy = orgExist.CreatedBy
	return &org, nil
}

func (o *OrganizationService) FindOrganizationByID(orgID uint, userID uint) (*dtos.Organization, *common.ErrorCodeMessage) {
	isAdmin, _, err := o.CheckUserRoleInOrganization(orgID, userID)
	if !isAdmin && err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidRequest,
			Message:     err.Error(),
		}
	}
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	orgDTO := mapper.OrganizationModelToDTO(org)
	return orgDTO, nil
}

func (o *OrganizationService) AddPeopleToOrganization(ctx context.Context, orgID uint, userID uint, req *request.AddPeopleToOrganizationRequest) ([]string, *common.ErrorCodeMessage) {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	_, isManager, _ := o.CheckUserRoleInOrganization(orgID, userID)
	if !isManager {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	emails := make([]*string, 0)
	mapEmailDept := make(map[string]uint, 0)
	for _, u := range req.Users {
		emails = append(emails, &u.Email)
		mapEmailDept[u.Email] = u.DepartmentID
	}
	validUsers, err := o.userRepository.FindUsersNotInOrganization(emails)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if validUsers == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidUser,
			Message:     common.ErrMessageInvalidRequest,
		}
	}
	validEmails := make([]string, 0)
	for _, u := range validUsers {
		validEmails = append(validEmails, u.Email)
		if err := utils.SendOrganizationInvitation(o.emailConfig.Domain, orgID, org.Name, o.emailConfig.SenderEmail, o.emailConfig.SenderPassword, mapEmailDept[u.Email], u.Email); err != nil {
			logger.Error(ctx, err.Error())
		}
	}
	return validEmails, nil
}

func (o *OrganizationService) RemoveUserFromOrganization(ctx context.Context, orgID uint, userID uint, req *request.RemoveUserFromOrganizationRequest) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	_, isManager, _ := o.CheckUserRoleInOrganization(orgID, userID)
	if !isManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	if err := o.userRepository.RemoveUserFromOrganization(req.UserID); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (o *OrganizationService) AssignPeopleTobeManager(ctx context.Context, orgID uint, userID uint, req *request.AssignPeopleToManagerRequest) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	_, isManager, _ := o.CheckUserRoleInOrganization(orgID, userID)
	if !isManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	validUser, err := o.userRepository.FindUserInOrganization(req.Email, orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if validUser == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserNotInOrganization,
			Message:     common.ErrMessageUserNotInOrganization,
		}
	}
	if validUser.Status != 1 {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserIsInactive,
			Message:     common.ErrMessageUserIsInactive,
		}
	}
	if err := o.userRepository.UpdateUserRoleManager(validUser.ID, true); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	// TODO: send email or send noti to assigned people
	return nil
}

func (o *OrganizationService) RecallPeopleTobeManager(ctx context.Context, orgID uint, userID uint, req *request.RecallPeopleManagerRequest) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	user, err := o.userRepository.FinduserByID(userID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.Email != org.CreatedBy {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	validUser, err := o.userRepository.FindUserInOrganization(req.Email, orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if validUser == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserNotInOrganization,
			Message:     common.ErrMessageUserNotInOrganization,
		}
	}
	if validUser.Status != 1 {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserIsInactive,
			Message:     common.ErrMessageUserIsInactive,
		}
	}
	if err := o.userRepository.UpdateUserRoleManager(validUser.ID, false); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	// TODO: send email or send noti to assigned people
	return nil
}

func (o *OrganizationService) AcceptOrganizationInvitation(orgID uint, deptID uint, userEmail string) *common.ErrorCodeMessage {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if org == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	user, err := o.userRepository.FindUserByEmail(userEmail)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserNotFound,
			Message:     common.ErrMessageInvalidUser,
		}
	}
	if user.OrganizationID != 0 {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyInOtherOrganization,
			Message:     common.ErrMessageUserAlreadyInOtherOrganization,
		}
	}
	if err := o.userRepository.AddPeopleOrganization(user.ID, orgID, deptID); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (o *OrganizationService) CheckUserRoleInOrganization(orgID uint, userID uint) (bool, bool, error) {
	user, err := o.userRepository.FinduserByID(userID)
	if err != nil {
		return false, false, err
	}
	if user.IsAdmin {
		if user.OrganizationID == 0 || user.OrganizationID != orgID {
			return true, false, errors.New(common.ErrMessageCannotAccessToOrganization)
		}
		if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsOrganizationManager {
			return true, true, nil
		}
	} else {
		if user.OrganizationID == 0 || user.OrganizationID != orgID {
			return false, false, errors.New(common.ErrMessageCannotAccessToOrganization)
		}
		if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsOrganizationManager {
			return false, true, nil
		}
	}

	return false, false, nil
}
