package services

import (
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos/request"
	"errors"
	"time"
)

type IOrganizationService interface {
	CreateOrganization(req *request.CreateOrganizationRequest) error
}

type OrganizationService struct {
	organizationRepository postgres.IOrganizationRepository
}

func NewOrganizationService(orgRepo postgres.IOrganizationRepository) IOrganizationService {
	return &OrganizationService{organizationRepository: orgRepo}
}

func (o *OrganizationService) CreateOrganization(req *request.CreateOrganizationRequest) error {
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
	return nil
}
