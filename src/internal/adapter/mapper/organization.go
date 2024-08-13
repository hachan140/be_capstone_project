package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func OrganizationModelToDTO(orgModel *model.Organization) *dtos.Organization {
	percent := orgModel.DataUsed * 100 / orgModel.LimitData
	return &dtos.Organization{
		ID:              orgModel.ID,
		Name:            orgModel.Name,
		Description:     orgModel.Description,
		Status:          orgModel.Status,
		LimitData:       orgModel.LimitData,
		DataUsed:        orgModel.DataUsed,
		LimitToken:      orgModel.LimitToken,
		TokenUsed:       orgModel.TokenUsed,
		PercentDataUsed: percent,
		IsOpenai:        orgModel.IsOpenai,
		CreatedAt:       orgModel.CreatedAt,
		CreatedBy:       orgModel.CreatedBy,
		UpdatedAt:       orgModel.UpdatedAt,
		UpdatedBy:       orgModel.UpdatedBy,
	}
}
