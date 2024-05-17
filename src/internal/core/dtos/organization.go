package dtos

import "time"

type Organization struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}
