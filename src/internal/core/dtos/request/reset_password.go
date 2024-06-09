package request

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/utils"
	"net/http"
)

type ResetPasswordRequest struct {
	Email string `json:"email,omitempty"`
}

func (r *ResetPasswordRequest) Validate() *common.ErrorCodeMessage {
	if r.Email == "" {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidRequest,
			Message:     common.ErrMessageInvalidEmail,
		}
	}
	return nil
}

type ResetPassword struct {
	NewPassword string `json:"new_password,omitempty"`
}

func (r *ResetPassword) Validate() *common.ErrorCodeMessage {
	isValidPassword := utils.ValidatePassword(r.NewPassword)
	if !isValidPassword {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeInvalidPassword,
			Message:     common.ErrMessageInvalidPassword,
		}
	}
	return nil
}
