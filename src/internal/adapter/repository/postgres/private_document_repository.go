package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type IPrivateDocumentRepository interface {
	GetPrivateDocument(userID uint, docID uint) (*model.PrivateDocs, error)
}

type PrivateDocumentRepository struct {
	storage *storage.Database
}

func NewPrivateDocumentRepository(storage *storage.Database) IPrivateDocumentRepository {
	return &PrivateDocumentRepository{storage: storage}
}

func (p *PrivateDocumentRepository) GetPrivateDocument(userID uint, docID uint) (*model.PrivateDocs, error) {
	var privateDoc *model.PrivateDocs
	err := p.storage.Raw("select * from private_docs where user_id = ? and doc_id = ? and status = 1", userID, docID).Scan(&privateDoc).Error
	if err != nil {
		return nil, err
	}
	return privateDoc, nil
}
