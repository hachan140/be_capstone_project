package request

import (
	"be-capstone-project/src/internal/core/common"
	"errors"
)

type CreateOrganizationRequest struct {
	Name        string
	Description string
	CreatedBy   string
}

func (o *CreateOrganizationRequest) Validate() error {
	if o.Name == "" {
		return errors.New(common.ErrMessageInvalidOrganizationName)
	}
	return nil
}

type UpdateOrganizationRequest struct {
	Name        *string
	Description *string
	Status      *int
	UpdatedBy   string
}

func (o *UpdateOrganizationRequest) Validate() error {
	if *o.Name == "" {
		return errors.New(common.ErrMessageInvalidOrganizationName)
	}
	return nil
}

type AddPeopleToOrganizationRequest struct {
	Users []*UserInfo `json:"users"`
}

type UserInfo struct {
	Email        string `json:"email"`
	DepartmentID uint   `json:"department_id"`
}
