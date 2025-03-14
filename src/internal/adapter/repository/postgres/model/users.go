package model

import "time"

type User struct {
	ID                    uint
	FirstName             string
	LastName              string
	Email                 string
	Password              string
	Gender                bool
	Status                int
	IsAdmin               bool
	IsOrganizationManager bool
	IsDeptManager         bool
	OrganizationID        uint
	DeptID                uint
	IsSocial              bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
