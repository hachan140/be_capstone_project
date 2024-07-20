package model

import "time"

type Product struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Price       float64
	Quantity    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
