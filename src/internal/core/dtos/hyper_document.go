package dtos

import "time"

type HyperDocument struct {
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
	Duration    int64
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
}
