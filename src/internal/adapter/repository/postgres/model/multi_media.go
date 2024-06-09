package model

import "time"

type MultiMedia struct {
	ID          uint
	Title       string
	Description string
	FilePath    string
	Duration    int64
	Status      int
	CategoryID  uint
	Type        string
	FileID      string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
}
