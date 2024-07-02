package dtos

import "time"

type Document struct {
	ID              uint
	Title           string
	Description     string
	Content         string
	CategoryID      uint
	TotalPage       int
	Status          int
	Type            string
	Duration        int64
	FilePath        string
	FileID          string
	StorageCapacity int64
	StorageUnit     string
	AccessType      int
	DeptID          uint
	OrganizationID  uint
	CreatedAt       time.Time
	CreatedBy       string
	UpdatedAt       time.Time
}
