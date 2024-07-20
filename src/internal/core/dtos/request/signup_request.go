package request

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/utils"
	"errors"
)

type SignUpRequest struct {
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   bool   `json:"gender"`
}

func (s *SignUpRequest) Validate() error {
	if s.FistName == "" || s.LastName == "" {
		return errors.New(common.ErrMessageInvalidName)
	}
	isValidEmail := utils.ValidateEmail(s.Email)
	if !isValidEmail {
		return errors.New(common.ErrMessageInvalidEmail)
	}
	isValidPassword := utils.ValidatePassword(s.Password)
	if !isValidPassword {
		return errors.New(common.ErrMessageInvalidPassword)
	}
	return nil
}

type VerifyEmail struct {
	Email string `form:"email"`
}
