package model

import "time"

type Document struct {
	ID          uint
	Title       string
	Description string
	Content     string
	CategoryID  uint
	TotalPage   int
	Status      int
	Type        string
	FilePath    string
	FileID      string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
}
