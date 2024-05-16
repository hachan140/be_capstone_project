package request

import (
	"be-capstone-project/src/internal/core/common"
	"errors"
)

type CreateOrganizationRequest struct {
	Name        string
	Description string
	CreatedBy   string
}

func (o *CreateOrganizationRequest) Validate() error {
	if o.Name == "" {
		return errors.New(common.ErrMessageInvalidOrganizationName)
	}
	return nil
}
