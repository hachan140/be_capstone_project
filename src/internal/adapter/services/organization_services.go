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
	CreateOrganization(userId uint, req *request.CreateOrganizationRequest) error
	UpdateOrganization(orgID uint, userID uint, req *request.UpdateOrganizationRequest) error
	FindOrganizationByID(orgID uint, userID uint) (*dtos.Organization, error)
	CheckUserRoleInOrganization(orgID uint, userID uint) (bool, error)
	AddPeopleToOrganization(ctx context.Context, orgID uint, userID uint, emails []*string) ([]string, *common.ErrorCodeMessage)
	AcceptOrganizationInvitation(orgID uint, userEmail string) *common.ErrorCodeMessage
}

type OrganizationService struct {
	organizationRepository postgres.IOrganizationRepository
	userRepository         postgres.IUserRepository
	emailConfig            common_configs.EmailConfig
}

func NewOrganizationService(orgRepo postgres.IOrganizationRepository, userRepository postgres.IUserRepository, emailConfig common_configs.EmailConfig) IOrganizationService {
	return &OrganizationService{
		organizationRepository: orgRepo,
		userRepository:         userRepository,
		emailConfig:            emailConfig,
	}
}

func (o *OrganizationService) CreateOrganization(userId uint, req *request.CreateOrganizationRequest) error {
	user, err := o.userRepository.FinduserByID(userId)
	if err != nil {
		return err
	}
	if user.OrganizationID != 0 {
		return errors.New(common.ErrMessageUserAlreadyInOtherOrganization)
	}
	orgByName, err := o.organizationRepository.FindOrganizationByName(req.Name)
	if err != nil {
		return err
	}
	if orgByName != nil {
		return errors.New(common.ErrMessageOrganizationExisted)
	}
	orgModel := &model.Organization{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   req.CreatedBy,
		Status:      3,
	}
	if err := o.organizationRepository.CreateOrganization(orgModel); err != nil {
		return err
	}
	/*	err = o.userRepository.UpdateUserOrganizationRole(userId, orgModel.ID, true)
		if err != nil {
			return err
		}*/
	return nil
}

func (o *OrganizationService) UpdateOrganization(orgID uint, userID uint, req *request.UpdateOrganizationRequest) error {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return err
	}
	if org == nil {
		return errors.New(common.ErrMessageOrganizationNotExist)
	}
	isOrganizationManager, err := o.CheckUserRoleInOrganization(orgID, userID)
	if !isOrganizationManager || err != nil {
		return errors.New(common.ErrMessageCannotAccessToOrganization)
	}
	if req.Name != nil && *req.Name != org.Name {
		orgByName, err := o.organizationRepository.FindOrganizationByName(*req.Name)
		if err != nil {
			return err
		}
		if orgByName != nil {
			return errors.New(common.ErrMessageOrganizationExisted)
		}
	}

	orgToUpdate, err := o.buildOrganizationUpdateQuery(org, req)
	if err != nil {
		return err
	}
	if err := o.organizationRepository.UpdateOrganization(orgID, orgToUpdate); err != nil {
		return err
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

func (o *OrganizationService) FindOrganizationByID(orgID uint, userID uint) (*dtos.Organization, error) {
	_, err := o.CheckUserRoleInOrganization(orgID, userID)
	if err != nil {
		return nil, err
	}
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return nil, errors.New(common.ErrMessageInvalidRequest)
	}
	if org == nil {
		return nil, errors.New(common.ErrMessageOrganizationNotExist)
	}
	orgDTO := mapper.OrganizationModelToDTO(org)
	return orgDTO, nil
}

func (o *OrganizationService) AddPeopleToOrganization(ctx context.Context, orgID uint, userID uint, emails []*string) ([]string, *common.ErrorCodeMessage) {
	org, err := o.organizationRepository.FindOrganizationByID(orgID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	if org == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeOrganizationNotExist,
			Message:     common.ErrMessageOrganizationNotExist,
		}
	}
	isManager, _ := o.CheckUserRoleInOrganization(orgID, userID)
	if !isManager {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	validEmails, err := o.userRepository.FindUsersNotInOrganization(emails)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	for _, email := range validEmails {
		if err := utils.SendOrganizationInvitation(orgID, org.Name, o.emailConfig.SenderEmail, o.emailConfig.SenderPassword, email); err != nil {
			logger.Error(ctx, err.Error())
		}
	}
	return validEmails, nil
}

func (o *OrganizationService) AcceptOrganizationInvitation(orgID uint, userEmail string) *common.ErrorCodeMessage {
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
	if err := o.userRepository.AddPeopleOrganization(user.ID, orgID); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (o *OrganizationService) CheckUserRoleInOrganization(orgID uint, userID uint) (bool, error) {
	user, err := o.userRepository.FinduserByID(userID)
	if err != nil {
		return false, err
	}
	if user.OrganizationID == 0 || user.OrganizationID != orgID {
		return false, errors.New(common.ErrMessageCannotAccessToOrganization)
	}
	if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsOrganizationManager {
		return true, nil
	}
	return false, nil
}

func (o *OrganizationService) SendOrganizationInvitationToUsers(senderEmail string, senderPassword string, receiverEmail []string, orgName string) error {

	return nil
}
