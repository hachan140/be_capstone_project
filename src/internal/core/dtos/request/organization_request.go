package request

import (
	"be-capstone-project/src/internal/core/common"
	"errors"
)

type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	IsOpenai    bool   `json:"is_openai"`
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

type UpdateOrganizationStatusRequest struct {
	Status    *int
	UpdatedBy string
}

type AddPeopleToOrganizationRequest struct {
	Users []*UserInfo `json:"users"`
}

type RemoveUserFromOrganizationRequest struct {
	UserID uint `json:"user_id"`
}

func (r *RemoveUserFromOrganizationRequest) Validate() error {
	if r.UserID == 0 {
		return errors.New("Invalid user_id")
	}
}

type UserInfo struct {
	Email        string `json:"email"`
	DepartmentID uint   `json:"department_id"`
}
