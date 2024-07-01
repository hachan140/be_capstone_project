package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func CategoryModelToDTO(model *model.Category) *dtos.Category {
	if model == nil {
		return nil
	}
	return &dtos.Category{
		ID:               model.ID,
		Name:             model.Name,
		Description:      model.Description,
		ParentCategoryID: model.ParentCategoryID,
		OrganizationID:   model.OrganizationID,
		DepartmentID:     model.DepartmentID,
		Status:           model.Status,
		CreatedBy:        model.CreatedBy,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}
