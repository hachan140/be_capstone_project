package model

import "time"

type Organization struct {
	ID          uint
	Name        string
	Description string
	Status      int
	LimitData   int64
	DataUsed    int64
	IsOpenai    bool
	LimitToken  int64
	TokenUsed   int64
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}
