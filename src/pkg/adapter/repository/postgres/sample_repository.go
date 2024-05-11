package postgres

import (
	"be-capstone-project/src/pkg/adapter/repository/postgres/model"
	"be-capstone-project/src/pkg/core/storage"
	"context"
)

type ISampleRepository interface {
	CreateSampleRepository(ctx context.Context, sample *model.Sample) error
}

type SampleRepository struct {
	storage *storage.Database
}

func NewSampleRepository(db *storage.Database) ISampleRepository {
	return &SampleRepository{
		storage: db,
	}
}

func (r *SampleRepository) CreateSampleRepository(ctx context.Context, sample *model.Sample) error {
	result := r.storage.Create(&sample)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
