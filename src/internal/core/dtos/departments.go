package dtos

import "time"

type Department struct {
	ID             uint      `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Status         int       `json:"status,omitempty"`
	OrganizationID uint      `json:"organization_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by,omitempty"`
}
