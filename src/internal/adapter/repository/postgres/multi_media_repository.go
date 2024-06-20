package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/storage"
	"context"
)

type IMultiMediaRepository interface {
	FilterMultiMedia(ctx context.Context, query string, params []interface{}, req request.HyperDocumentFilterParam) ([]*model.MultiMedia, error)
}

type MultiMediaRepository struct {
	storage *storage.Database
}

func NewMultiMediaRepository(storage *storage.Database) IMultiMediaRepository {
	return &MultiMediaRepository{storage: storage}
}

func (d *MultiMediaRepository) FilterMultiMedia(ctx context.Context, query string, params []interface{}, req request.HyperDocumentFilterParam) ([]*model.MultiMedia, error) {
	var multimedia []*model.MultiMedia
	sqlQuery := `(SELECT distinct multimedia.* FROM multimedia`
	sqlQuery += query
	sqlQuery += ` ORDER BY multimedia.id DESC `
	sqlQuery += " LIMIT ?"
	params = append(params, req.PageSize)
	sqlQuery += " OFFSET ?)"
	params = append(params, (req.Page-1)*req.PageSize)
	if err := d.storage.WithContext(ctx).Raw(sqlQuery, params...).Find(&multimedia); err.Error != nil {
		return nil, err.Error
	}
	return multimedia, nil
}
