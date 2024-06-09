package postgres

import "be-capstone-project/src/internal/core/storage"

type IHyperDocumentRepository interface {
}

type HyperDocumentRepository struct {
	storage *storage.Database
}

func NewHyperDocumentRepository(storage *storage.Database) IHyperDocumentRepository {
	return &HyperDocumentRepository{storage: storage}
}

func (h *HyperDocumentRepository) FilterHyperDocument()
