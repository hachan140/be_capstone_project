package request

import (
	"errors"
	"strings"
)

type CreateSampleRequest struct {
	Name        string `json:"name" validate:"required,gte=0" message:"name is required"`
	StudentID   string `json:"student_id"`
	Email       string `json:"email" `
	PhoneNumber string `json:"phone_number"`
}

func (g *CreateSampleRequest) Validate() error {
	if g.Name == "" || g.Email == "" || g.PhoneNumber == "" {
		return errors.New("name, email and phone_number params are required")
	}
	if g.Email != "" {
		g.Email = strings.ToLower(g.Email)
	}
	return nil
}
