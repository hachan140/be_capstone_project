package request

import (
	"be-capstone-project/src/internal/core/common"
	"errors"
)

type CreateCategoryRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	ParentID       uint   `json:"parent_id"`
	OrganizationID uint   `json:"organization_id"`
	DepartmentID   uint   `json:"department_id"`
	CreatedBy      string `json:"created_by"`
}

func (c *CreateCategoryRequest) Validate() error {
	if c.Name == "" {
		return errors.New(common.ErrMessageInvalidCategoryName)
	}
	if c.OrganizationID == 0 {
		return errors.New("Invalid organization ID")
	}
	return nil
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *uint   `json:"parent_id"`
	Status      *int    `json:"status"`
	UpdatedBy   string  `json:"updated_by"`
}

func (c *UpdateCategoryRequest) Validate() error {
	if *c.Name == "" {
		return errors.New(common.ErrMessageInvalidCategoryName)
	}
	return nil
}

type UpdateCategoryStatusRequest struct {
	Status    *int   `json:"status"`
	UpdatedBy string `json:"updated_by"`
}

type UpdateDepartmentStatusRequest struct {
	Status    *int   `json:"status"`
	UpdatedBy string `json:"updated_by"`
}
