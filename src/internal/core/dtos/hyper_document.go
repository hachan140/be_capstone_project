package dtos

import "time"

type Document struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Content         string    `json:"content"`
	CategoryID      uint      `json:"category_id"`
	TotalPage       int       `json:"total_page"`
	Status          int       `json:"status"`
	Type            string    `json:"type"`
	Duration        int64     `json:"duration"`
	FilePath        string    `json:"file_path"`
	FileID          string    `json:"file_id"`
	StorageCapacity int64     `json:"storage_capacity"`
	StorageUnit     string    `json:"storage_unit"`
	AccessType      int       `json:"access_type"`
	DeptID          uint      `json:"dept_id"`
	OrganizationID  uint      `json:"organization_id"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}
