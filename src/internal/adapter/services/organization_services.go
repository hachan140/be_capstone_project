package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"errors"
	"time"
)

type IOrganizationService interface {
	CreateOrganization(userId uint, req *request.CreateOrganizationRequest) error
	UpdateOrganization(orgID uint, userID uint, req *request.UpdateOrganizationRequest) error
	FindOrganizationByID(orgID uint, userID uint) (*dtos.Organization, error)
}

type OrganizationService struct {
	organizationRepository postgres.IOrganizationRepository
	userRepository         postgres.IUserRepository
}

func NewOrganizationService(orgRepo postgres.IOrganizationRepository, userRepository postgres.IUserRepository) IOrganizationService {
	return &OrganizationService{
		organizationRepository: orgRepo,
		userRepository:         userRepository,
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
	}
	if err := o.organizationRepository.CreateOrganization(orgModel); err != nil {
		return err
	}
	err = o.userRepository.UpdateUserOrganizationRole(userId, orgModel.ID, true)
	if err != nil {
		return err
	}
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
	isOrganizationManager, err := o.checkUserRoleInOrganization(orgID, userID)
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
	_, err := o.checkUserRoleInOrganization(orgID, userID)
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

func (o *OrganizationService) checkUserRoleInOrganization(orgID uint, userID uint) (bool, error) {
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
