package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type IOrganizationRepository interface {
	CreateOrganization(org *model.Organization) error
	FindOrganizationByID(ID uint) (*model.Organization, error)
	FindOrganizationByName(name string) (*model.Organization, error)
	UpdateOrganization(orgID uint, org *model.Organization) error
}

type OrganizationRepository struct {
	storage *storage.Database
}

func NewOrganizationRepository(storage *storage.Database) IOrganizationRepository {
	return &OrganizationRepository{
		storage: storage,
	}
}

func (o *OrganizationRepository) CreateOrganization(org *model.Organization) error {
	err := o.storage.Model(org).Create(&org).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrganizationRepository) FindOrganizationByID(ID uint) (*model.Organization, error) {
	var org *model.Organization
	err := o.storage.Raw("select * from organizations where id = ?", ID).Scan(&org).Error
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (o *OrganizationRepository) FindOrganizationByName(name string) (*model.Organization, error) {
	var org *model.Organization
	err := o.storage.Raw("select * from organizations where name = ?", name).Scan(&org).Error
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (o *OrganizationRepository) UpdateOrganization(orgID uint, org *model.Organization) error {
	err := o.storage.Model(org).Where("id = ?", orgID).Updates(org).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrganizationRepository) DeleteOrganizationByID(ID uint) error {
	err := o.storage.Raw("delete from organizations where id = ?", ID).Error
	if err != nil {
		return err
	}
	return nil
}
