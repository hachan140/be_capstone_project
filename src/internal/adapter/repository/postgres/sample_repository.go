package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
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
	_, err := r.storage.NewInsert().Model(sample).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
