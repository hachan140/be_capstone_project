package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/storage"
	"context"
)

type IDocumentRepository interface {
	FilterDocument(ctx context.Context, query string, params []interface{}, req request.HyperDocumentFilterParam) ([]*model.Document, error)
}

type DocumentRepository struct {
	storage *storage.Database
}

func NewDocumentRepository(storage *storage.Database) IDocumentRepository {
	return &DocumentRepository{storage: storage}
}

func (d *DocumentRepository) FilterDocument(ctx context.Context, query string, params []interface{}, req request.HyperDocumentFilterParam) ([]*model.Document, error) {
	var documents []*model.Document
	sqlQuery := `(SELECT distinct documents.* FROM documents`
	sqlQuery += query
	sqlQuery += ` ORDER BY documents.id DESC `
	sqlQuery += " LIMIT ?"
	params = append(params, req.PageSize)
	sqlQuery += " OFFSET ?)"
	params = append(params, (req.Page-1)*req.PageSize)
	if err := d.storage.WithContext(ctx).Raw(sqlQuery, params...).Find(&documents); err.Error != nil {
		return nil, err.Error
	}
	return documents, nil
}
