package dtos

import "time"

type Organization struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          int       `json:"status"`
	IsOpenai        bool      `json:"is_openai"`
	LimitData       int64     `json:"limit_data"`
	DataUsed        int64     `json:"data_used"`
	LimitToken      int64     `json:"limit_token"`
	TokenUsed       int64     `json:"token_used"`
	PercentDataUsed int64     `json:"percent_data_used"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
}
