package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func SampleModelToDTO(model *model.Sample) *dtos.Sample {
	if model == nil {
		return &dtos.Sample{}
	}
	return &dtos.Sample{
		ID:          model.ID,
		Name:        model.Name,
		StudentID:   model.StudentID,
		Email:       model.Email,
		PhoneNumber: model.PhoneNumber,
	}
}
