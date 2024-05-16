package request

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/utils"
	"errors"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {

	isValidEmail := utils.ValidateEmail(l.Email)
	if !isValidEmail {
		return errors.New(common.ErrMessageInvalidEmail)
	}

	return nil
}
