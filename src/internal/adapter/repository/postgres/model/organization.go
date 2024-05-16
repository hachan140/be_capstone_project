package model

import "time"

type Organization struct {
	ID          uint
	Name        string
	Description string
	Status      int
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}
