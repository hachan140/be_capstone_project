package request

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/utils"
	"errors"
)

type SocialLoginRequest struct {
	Email     string `json:"email"`
	IsSocial  bool   `json:"is_social"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *SocialLoginRequest) Validate() error {
	isValidEmail := utils.ValidateEmail(s.Email)
	if !isValidEmail {
		return errors.New(common.ErrMessageInvalidEmail)
	}
	if !s.IsSocial {
		return errors.New(common.ErrMessageInvalidRequest)
	}
	return nil
}
