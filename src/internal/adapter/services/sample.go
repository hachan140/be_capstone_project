package services

import (
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
)

type ISampleService interface {
	CreateSampleService(ctx context.Context, request *request.CreateSampleRequest) error
}

type SampleService struct {
	sampleRepo postgres.ISampleRepository
}

func NewSampleService(sampleRepo postgres.ISampleRepository) ISampleService {
	return &SampleService{sampleRepo: sampleRepo}
}

func (s *SampleService) CreateSampleService(ctx context.Context, request *request.CreateSampleRequest) error {
	sample := &model.Sample{
		Name:        request.Name,
		StudentID:   request.StudentID,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	}
	if err := s.sampleRepo.CreateSampleRepository(ctx, sample); err != nil {
		return err
	}
	return nil
}
