package model

import "time"

type Category struct {
	ID               uint
	Name             string
	Description      string
	ParentCategoryID uint
	OrganizationID   uint
	DepartmentID     uint
	Status           int
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
