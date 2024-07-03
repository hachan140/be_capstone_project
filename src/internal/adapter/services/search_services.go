package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"be-capstone-project/src/internal/core/dtos/response"
	"net/http"
	"strings"
)

type ISearchService interface {
	SearchDocumentAndOrNot(req *request.SearchAndOrNotRequest, userID uint) ([]*dtos.Document, *common.ErrorCodeMessage)
	GetSearchKeywords(userID uint, req *request.SearchHistoryRequest) (*response.SearchHistoryResponse, *common.ErrorCodeMessage)
}

type SearchService struct {
	privateDocRepo          postgres.IPrivateDocumentRepository
	documentRepo            postgres.IDocumentRepository
	userRepository          postgres.IUserRepository
	searchHistoryRepository postgres.ISearchHistoryRepository
}

func NewSearchService(privateDoc postgres.IPrivateDocumentRepository, documentRepo postgres.IDocumentRepository, userRepo postgres.IUserRepository,
	searchHistoryRepo postgres.ISearchHistoryRepository) ISearchService {
	return &SearchService{
		privateDocRepo:          privateDoc,
		documentRepo:            documentRepo,
		userRepository:          userRepo,
		searchHistoryRepository: searchHistoryRepo,
	}
}

func (s *SearchService) SearchDocumentAndOrNot(req *request.SearchAndOrNotRequest, userID uint) ([]*dtos.Document, *common.ErrorCodeMessage) {
	user, err := s.userRepository.FinduserByID(userID)
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

	titles, err := s.documentRepo.SearchDocumentTitles()
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	var results []string
	for _, title := range titles {
		if s.containsAll(title, req.AndKeyWords) && s.containsAny(title, req.OrKeyWords) && s.containsNone(title, req.NotKeywords) {
			results = append(results, title)
		}
	}
	documents, err := s.documentRepo.GetDocumentByTitles(results)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if documents == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeDocumentNotFound,
			Message:     common.ErrMessageDocumentNotFound,
		}
	}
	documentRes := s.filterDocumentAccessType(userID, user.OrganizationID, user.DeptID, documents)
	documentDtos := mapper.DocumentsToHyperDocumentDTOs(documentRes)
	return documentDtos, nil
}

func (s *SearchService) containsAll(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(text, keyword) {
			return false
		}
	}
	return true
}

func (s *SearchService) containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (s *SearchService) containsNone(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return false
		}
	}
	return true
}

// 1: public
// 2: inside org
// 3: inside department
// 4: shared people

func (s *SearchService) filterDocumentAccessType(userID uint, orgUserID uint, userDeptID uint, documents []*model.Document) []*model.Document {

	docRes := make([]*model.Document, 0)
	for _, d := range documents {
		if d.AccessType == 1 {
			docRes = append(docRes, d)
		}
		if d.OrganizationID == orgUserID {
			switch d.AccessType {
			case 2:
				if d.OrganizationID != 0 {
					docRes = append(docRes, d)
				}

			case 3:
				if d.OrganizationID != 0 && d.DeptID != 0 && d.DeptID == userDeptID {
					docRes = append(docRes, d)
				}

			case 4:
				doc, err := s.privateDocRepo.GetPrivateDocument(userID, d.ID)
				if err != nil {
					return nil
				}
				if doc != nil {
					docRes = append(docRes, d)
				}
			default:
				return nil
			}
		}
		continue
	}
	return docRes
}

func (s *SearchService) GetSearchKeywords(userID uint, req *request.SearchHistoryRequest) (*response.SearchHistoryResponse, *common.ErrorCodeMessage) {
	keywords, err := s.searchHistoryRepository.GetAllSearchHistoryPersonalize(userID, req.Keyword)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	res := &response.SearchHistoryResponse{Keywords: keywords}
	return res, nil
}
