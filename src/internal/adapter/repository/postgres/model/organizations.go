package model

import "time"

type Organization struct {
	ID          uint
	Name        string
	Description string
	Status      int
	LimitData   int64
	DataUsed    int64
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}
