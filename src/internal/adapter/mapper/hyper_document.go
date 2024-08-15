package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func DocumentToHyperDocumentDTO(document *model.Document) *dtos.Document {
	return &dtos.Document{
		ID:              document.ID,
		Title:           document.Title,
		Description:     document.Description,
		Content:         document.Content,
		AccessType:      document.AccessType,
		CategoryID:      document.CategoryID,
		TotalPage:       document.TotalPage,
		Status:          document.Status,
		Type:            document.Type,
		FilePath:        document.FilePath,
		FileID:          document.FileID,
		Duration:        0,
		StorageCapacity: document.StorageCapacity,
		StorageUnit:     document.StorageUnit,
		CreatedAt:       document.CreatedAt,
		CreatedBy:       document.CreatedBy,
		UpdatedAt:       document.UpdatedAt,
	}
}

func DocumentsToHyperDocumentDTOs(document []*model.Document) []*dtos.Document {
	if document == nil {
		return nil
	}
	documentRes := make([]*dtos.Document, 0)
	for _, document := range document {
		dRes := &dtos.Document{
			ID:          document.ID,
			Title:       document.Title,
			Description: document.Description,
			Content:     document.Content,
			CategoryID:  document.CategoryID,
			TotalPage:   document.TotalPage,
			Status:      document.Status,
			Type:        document.Type,
			FilePath:    document.FilePath,
			FileID:      document.FileID,
			Duration:    0,
			CreatedAt:   document.CreatedAt,
			CreatedBy:   document.CreatedBy,
			UpdatedAt:   document.UpdatedAt,
		}
		documentRes = append(documentRes, dRes)
	}
	return documentRes
}

func MultimediaToHyperDocumentDTO(multimedia *model.MultiMedia) *dtos.Document {
	return &dtos.Document{
		ID:          multimedia.ID,
		Title:       multimedia.Title,
		Description: multimedia.Description,
		Content:     "",
		CategoryID:  multimedia.CategoryID,
		TotalPage:   0,
		Status:      multimedia.Status,
		Type:        multimedia.Type,
		FilePath:    multimedia.FilePath,
		FileID:      multimedia.FileID,
		Duration:    multimedia.Duration,
		CreatedAt:   multimedia.CreatedAt,
		CreatedBy:   multimedia.CreatedBy,
		UpdatedAt:   multimedia.UpdatedAt,
	}
}
