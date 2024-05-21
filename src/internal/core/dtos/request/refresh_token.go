package request

import (
	"be-capstone-project/src/internal/core/common"
	"errors"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New(common.ErrMessageInvalidRequest)
	}
	return nil
}
