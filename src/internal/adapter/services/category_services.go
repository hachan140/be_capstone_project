package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	model2 "be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
	"net/http"
	"time"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, userID uint, req *request.CreateCategoryRequest) *common.ErrorCodeMessage
	ListCategories(ctx context.Context, orgID uint, userID uint, req *request.GetListCategoryRequest) ([]*dtos.Category, *common.ErrorCodeMessage)
	GetCategoryByID(ctx context.Context, id uint, userID uint) (*dtos.Category, *common.ErrorCodeMessage)
	UpdateCategoryByID(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryRequest) *common.ErrorCodeMessage
}

type CategoryService struct {
	categoryRepo   postgres.ICategoryRepository
	userRepository postgres.IUserRepository
}

func NewCategoryService(categoryRepo postgres.ICategoryRepository, userRepo postgres.IUserRepository) ICategoryService {
	return &CategoryService{categoryRepo: categoryRepo, userRepository: userRepo}
}

func (c *CategoryService) CreateCategory(ctx context.Context, userID uint, req *request.CreateCategoryRequest) *common.ErrorCodeMessage {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	isOrgManager, errC := c.CheckUserRoleInOrganization(req.OrganizationID, user.ID)
	if errC != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     errC.Message,
		}
	}
	if !isOrgManager || errC != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	catExisted, err := c.categoryRepo.FindCategoryByName(req.Name)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if catExisted != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCategoryExisted,
			Message:     common.ErrMessageCategoryExisted,
		}
	}
	model := &model2.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryID: req.ParentID,
		OrganizationID:   req.OrganizationID,
		DepartmentID:     req.DepartmentID,
		Status:           1,
		CreatedBy:        req.CreatedBy,
		CreatedAt:        time.Now(),
	}
	if err := c.categoryRepo.CreateCategory(model); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (c *CategoryService) ListCategories(ctx context.Context, orgID uint, userID uint, req *request.GetListCategoryRequest) ([]*dtos.Category, *common.ErrorCodeMessage) {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID != 0 && user.OrganizationID != orgID {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeUserAlreadyInOtherOrganization,
			Message:     common.ErrMessageUserAlreadyInOtherOrganization,
		}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	categories, err := c.categoryRepo.ListCategoryByOrganization(orgID, req.PageSize, offset)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	var categoryRes []*dtos.Category
	for _, c := range categories {
		cRes := &dtos.Category{
			ID:               c.ID,
			Name:             c.Name,
			Description:      c.Description,
			ParentCategoryID: c.ParentCategoryID,
			OrganizationID:   c.OrganizationID,
			DepartmentID:     c.DepartmentID,
			Status:           c.Status,
			CreatedBy:        c.CreatedBy,
			CreatedAt:        time.Now(),
			UpdatedAt:        c.UpdatedAt,
		}
		categoryRes = append(categoryRes, cRes)
	}
	return categoryRes, nil
}

func (c *CategoryService) GetCategoryByID(ctx context.Context, id uint, userID uint) (*dtos.Category, *common.ErrorCodeMessage) {

	cat, err := c.categoryRepo.FindCategoryByID(id)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if cat == nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCategoryNotFound,
			Message:     common.ErrMessageCategoryNotFound,
		}
	}
	_, errC := c.CheckUserRoleInOrganization(cat.OrganizationID, userID)
	if errC != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return mapper.CategoryModelToDTO(cat), nil
}

func (c *CategoryService) UpdateCategoryByID(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryRequest) *common.ErrorCodeMessage {
	category, err := c.categoryRepo.FindCategoryByID(catID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if category == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCategoryNotFound,
			Message:     common.ErrMessageCategoryNotFound,
		}
	}
	isOrgManager, errC := c.CheckUserRoleInOrganization(category.OrganizationID, userID)
	if !isOrgManager || errC != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if req.Name != nil {
		categoryByName, err := c.categoryRepo.FindCategoryByName(*req.Name)
		if err != nil {
			return &common.ErrorCodeMessage{
				HTTPCode:    http.StatusInternalServerError,
				ServiceCode: common.ErrCodeInternalError,
				Message:     err.Error(),
			}
		}
		if categoryByName != nil && categoryByName.Name != category.Name {
			return &common.ErrorCodeMessage{
				HTTPCode:    http.StatusBadRequest,
				ServiceCode: common.ErrCodeCategoryExisted,
				Message:     common.ErrMessageCategoryExisted,
			}
		}
	}

	categoryToUpdate := c.buildUpdateCategory(category, req)
	if err := c.categoryRepo.UpdateCategory(categoryToUpdate); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (c *CategoryService) CheckUserRoleInOrganization(orgID uint, userID uint) (bool, *common.ErrorCodeMessage) {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID == 0 || user.OrganizationID != orgID {
		return false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsOrganizationManager {
		return true, nil
	}
	return false, nil
}

func (c *CategoryService) buildUpdateCategory(existedCategory *model2.Category, req *request.UpdateCategoryRequest) *model2.Category {
	var category model2.Category
	if req == nil {
		return nil
	}
	if req.Name != nil {
		category.Name = *req.Name
	} else {
		category.Name = existedCategory.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	} else {
		category.Description = existedCategory.Description
	}
	if req.Status != nil {
		category.Status = *req.Status
	} else {
		category.Status = existedCategory.Status
	}
	if req.ParentID != nil {
		category.ParentCategoryID = *req.ParentID
	} else {
		category.ParentCategoryID = existedCategory.ParentCategoryID
	}
	category.ID = existedCategory.ID
	category.UpdatedAt = time.Now()
	return &category
}
