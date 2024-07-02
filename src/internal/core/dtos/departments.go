package dtos

import "time"

type Department struct {
	ID             uint
	Name           string
	Description    string
	Status         int
	OrganizationID uint
	CreatedAt      time.Time
	CreatedBy      string
}
