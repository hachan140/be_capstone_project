package model

import "time"

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Gender    bool
	Status    int
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
