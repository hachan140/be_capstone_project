package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func OrganizationModelToDTO(orgModel *model.Organization) *dtos.Organization {
	return &dtos.Organization{
		ID:          orgModel.ID,
		Name:        orgModel.Name,
		Description: orgModel.Description,
		Status:      orgModel.Status,
		IsOpenai:    orgModel.IsOpenai,
		CreatedAt:   orgModel.CreatedAt,
		CreatedBy:   orgModel.CreatedBy,
		UpdatedAt:   orgModel.UpdatedAt,
		UpdatedBy:   orgModel.UpdatedBy,
	}
}
