package model

import "time"

type SearchHistory struct {
	ID        uint
	UserID    uint
	Keywords  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      int
}
