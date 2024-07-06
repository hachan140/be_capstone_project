package dtos

import "time"

type Category struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ParentCategoryID uint      `json:"parent_category_id"`
	OrganizationID   uint      `json:"organization_id"`
	DepartmentID     uint      `json:"department_id"`
	Status           int       `json:"status"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
