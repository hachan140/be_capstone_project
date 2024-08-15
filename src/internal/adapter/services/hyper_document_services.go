package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
	"net/http"
	"strings"
)

type IHyperDocumentService interface {
	FilterHyperDocument(ctx context.Context, req request.HyperDocumentFilterParam, userID uint) ([]*dtos.Document, *common.ErrorCodeMessage)
}

type HyperDocumentService struct {
	documentRepository postgres.IDocumentRepository
	userRepository     postgres.IUserRepository
	privateDocRepo     postgres.IPrivateDocumentRepository
}

func NewHyperDocumentService(documentRepository postgres.IDocumentRepository, userRepository postgres.IUserRepository, privateDoc postgres.IPrivateDocumentRepository) IHyperDocumentService {
	return &HyperDocumentService{documentRepository: documentRepository, userRepository: userRepository, privateDocRepo: privateDoc}
}

func (h *HyperDocumentService) FilterHyperDocument(ctx context.Context, req request.HyperDocumentFilterParam, userID uint) ([]*dtos.Document, *common.ErrorCodeMessage) {
	user, err := h.userRepository.FinduserByID(userID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserNotFound,
			Message:     common.ErrMessageInvalidUser,
		}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 100 || req.PageSize <= 0 {
		req.PageSize = 100
	}
	queryDocuments, paramDocuments := h.BuildQueryFilterDocument(req)
	documents, err := h.documentRepository.FilterDocument(ctx, queryDocuments, paramDocuments, req)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	documents = h.filterDocumentAccessType(userID, user.OrganizationID, user.DeptID, documents)
	return mapper.DocumentsToHyperDocumentDTOs(documents), nil
}

func (h *HyperDocumentService) BuildQueryFilterDocument(req request.HyperDocumentFilterParam) (string, []interface{}) {
	query := ` WHERE 1 = 1 `
	var params []interface{}
	if req.Title != "" {
		query += ` AND LOWER(documents.title) like ?`
		params = append(params, "%"+strings.ToLower(req.Title)+"%")
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

func (h *HyperDocumentService) filterDocumentAccessType(userID uint, orgUserID uint, userDeptID uint, documents []*model.Document) []*model.Document {

	docRes := make([]*model.Document, 0)
	for _, d := range documents {
		if d.AccessType == 1 {
			docRes = append(docRes, d)
			continue
		}
		if d.OrganizationID == orgUserID {
			switch d.AccessType {
			case 2:
				if d.OrganizationID != 0 {
					docRes = append(docRes, d)
				}
				break
			case 3:
				if d.OrganizationID != 0 && d.DeptID != 0 && d.DeptID == userDeptID {
					docRes = append(docRes, d)
				}
				break
			case 4:
				doc, err := h.privateDocRepo.GetPrivateDocument(userID, d.ID)
				if err != nil {
					return nil
				}
				if doc != nil {
					docRes = append(docRes, d)
				}
				break
			default:
				return nil
			}
		}
		continue
	}
	return docRes
}
