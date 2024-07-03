package model

import "time"

type SearchHistory struct {
	ID        uint
	UserID    uint
	Keywords  uint
	CreatedAt time.Time
	Type      int
}
