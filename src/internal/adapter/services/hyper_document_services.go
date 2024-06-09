package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
	"net/http"
)

type IHyperDocumentService interface {
}

type HyperDocumentService struct {
	documentRepository   postgres.IDocumentRepository
	multimediaRepository postgres.IMultiMediaRepository
}

func NewHyperDocumentService(documentRepository postgres.IDocumentRepository, multimediaRepository postgres.IMultiMediaRepository) IHyperDocumentService {
	return &HyperDocumentService{documentRepository: documentRepository, multimediaRepository: multimediaRepository}
}

func (h *HyperDocumentService) FilterHyperDocument(ctx context.Context, req request.HyperDocumentFilterParam) ([]*dtos.HyperDocument, *common.ErrorCodeMessage) {
	queryDocuments, paramDocuments := h.BuildQueryFilterDocument(req)
	documents, err := h.documentRepository.FilterDocument(ctx, queryDocuments, paramDocuments, req)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	hyperDocuments := make([]*dtos.HyperDocument, 0)
	if documents != nil {
		for _, d := range documents {
			dDTO := mapper.DocumentToHyperDocumentDTO(d)
			hyperDocuments = append(hyperDocuments, dDTO)
		}
	}

	queryMultiMedias, paramMultiMedias := h.BuildQueryFilterMultiMedia(req)
	multiMedia, err := h.multimediaRepository.FilterMultiMedia(ctx, queryMultiMedias, paramMultiMedias, req)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if documents != nil {
		for _, d := range multiMedia {
			dDTO := mapper.MultimediaToHyperDocumentDTO(d)
			hyperDocuments = append(hyperDocuments, dDTO)
		}
	}
	return hyperDocuments, nil
}

func (h *HyperDocumentService) BuildQueryFilterDocument(req request.HyperDocumentFilterParam) (string, []interface{}) {
	query := `WHERE 1 = 1 `
	var params []interface{}
	if req.Title != "" {
		query += ` AND documents.title like '%?%'`
		params = append(params, req.Title)
	}
	if req.Type != "" {
		query += ` AND documents.type = ?`
		params = append(params, req.Type)
	}
	if req.CreatedBy != "" {
		query += ` AND documents.created_by = ?`
		params = append(params, req.CreatedBy)
	}
	if req.CreatedFromDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		query += " AND documents.created_at >= ?"
		params = append(params, req.CreatedFromDate)
	}
	if req.CreatedToDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		query += " AND documents.created_at <= ?"
		params = append(params, req.CreatedToDate)
	}
	return query, params
}

func (h *HyperDocumentService) BuildQueryFilterMultiMedia(req request.HyperDocumentFilterParam) (string, []interface{}) {
	query := `WHERE 1 = 1 `
	var params []interface{}
	if req.Title != "" {
		query += ` AND multimedia.title like '%?%'`
		params = append(params, req.Title)
	}
	if req.Type != "" {
		query += ` AND multimedia.type = ?`
		params = append(params, req.Type)
	}
	if req.CreatedBy != "" {
		query += ` AND multimedia.created_by = ?`
		params = append(params, req.CreatedBy)
	}
	if req.CreatedFromDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		query += " AND multimedia.created_at >= ?"
		params = append(params, req.CreatedFromDate)
	}
	if req.CreatedToDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		query += " AND multimedia.created_at <= ?"
		params = append(params, req.CreatedToDate)
	}
	return query, params
}
