package mapper

import (
	"be-capstone-project/src/pkg/adapter/repository/postgres/model"
	"be-capstone-project/src/pkg/core/dtos"
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
