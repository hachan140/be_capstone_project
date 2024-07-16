package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/storage"
	"context"
)

type IDocumentRepository interface {
	FilterDocument(ctx context.Context, query string, params []interface{}, req request.HyperDocumentFilterParam) ([]*model.Document, error)
	SearchDocumentTitles() ([]string, error)
	GetDocumentByTitles(titles []string) ([]*model.Document, error)
	UpdateDocumentStatusByOrganizationID(orgID uint, status int) error
	UpdateDocumentStatusByCategoryID(catID uint, status int) error
	UpdateDocumentStatusByDepartmentID(depID uint, status int) error
	UpdateDocument(document *model.Document) error
	GetDocumentsByCategoryID(catID uint) ([]*model.Document, error)
	GetDocumentsByDepartmentID(deptID uint) ([]*model.Document, error)
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

func (d *DocumentRepository) SearchDocumentTitles() ([]string, error) {
	var titles []string
	err := d.storage.Raw("select title from documents ").Scan(&titles).Error
	if err != nil {
		return nil, err
	}
	return titles, nil
}

func (d *DocumentRepository) GetDocumentByTitles(titles []string) ([]*model.Document, error) {
	var documents []*model.Document
	err := d.storage.Raw("select * from documents where title in ?", titles).Scan(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (d *DocumentRepository) GetDocumentsByCategoryID(catID uint) ([]*model.Document, error) {
	var documents []*model.Document
	err := d.storage.Raw("select * from documents where category_id = ? ", catID).Scan(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (d *DocumentRepository) GetDocumentsByDepartmentID(deptID uint) ([]*model.Document, error) {
	var documents []*model.Document
	err := d.storage.Raw("select * from documents where dept_id = ? ", deptID).Scan(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (d *DocumentRepository) UpdateDocumentStatusByOrganizationID(orgID uint, status int) error {
	err := d.storage.Exec("update documents set status = ? where organization_id = ?", status, orgID).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DocumentRepository) UpdateDocumentStatusByCategoryID(catID uint, status int) error {
	err := d.storage.Exec("update documents set status = ? where category_id = ?", status, catID).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DocumentRepository) UpdateDocumentStatusByDepartmentID(depID uint, status int) error {
	err := d.storage.Exec("update documents set status = ? where dept_ID = ?", status, depID).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DocumentRepository) UpdateDocument(document *model.Document) error {
	err := d.storage.Model(document).Where("id = ?", document.ID).Updates(document).Error
	if err != nil {
		return err
	}
	return nil
}
