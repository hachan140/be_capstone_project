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
	"strings"
	"time"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, userID uint, req *request.CreateCategoryRequest) *common.ErrorCodeMessage
	ListCategories(ctx context.Context, orgID uint, userID uint, req *request.GetListCategoryRequest) ([]*dtos.Category, *common.ErrorCodeMessage)
	GetCategoryByID(ctx context.Context, id uint, userID uint) (*dtos.Category, *common.ErrorCodeMessage)
	UpdateCategoryByID(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryRequest) *common.ErrorCodeMessage
	UpdateCategoryStatus(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryStatusRequest) *common.ErrorCodeMessage
	UpdateDepartmentStatus(ctx context.Context, userID uint, deptID uint, req *request.UpdateDepartmentStatusRequest) *common.ErrorCodeMessage
	SearchCategoryByName(ctx context.Context, name string, userID uint, deptID uint) ([]*dtos.Category, *common.ErrorCodeMessage)
}

type CategoryService struct {
	categoryRepo   postgres.ICategoryRepository
	userRepository postgres.IUserRepository
	documentRepo   postgres.IDocumentRepository
}

func NewCategoryService(categoryRepo postgres.ICategoryRepository, userRepo postgres.IUserRepository, documentRepo postgres.IDocumentRepository) ICategoryService {
	return &CategoryService{categoryRepo: categoryRepo, userRepository: userRepo, documentRepo: documentRepo}
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
	isDeptManager, isOrgManager, errC := c.CheckUserRoleInOrganization(req.OrganizationID, user.ID, req.DepartmentID)
	if errC != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     errC.Message,
		}
	}
	if !isDeptManager && !isOrgManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	catExisted, err := c.categoryRepo.FindCategoryByName(req.Name, req.DepartmentID)
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

func (c *CategoryService) ListCategories(ctx context.Context, depID uint, userID uint, req *request.GetListCategoryRequest) ([]*dtos.Category, *common.ErrorCodeMessage) {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	dept, err := c.categoryRepo.FindDepartmentByID(depID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID != dept.OrganizationID {
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
	categories, err := c.categoryRepo.ListCategoryByDepartment(depID, req.PageSize, offset)
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
			CreatedAt:        c.CreatedAt,
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
	_, _, errC := c.CheckUserRoleInOrganization(cat.OrganizationID, userID, 0)
	if errC != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return mapper.CategoryModelToDTO(cat), nil
}

func (c *CategoryService) SearchCategoryByName(ctx context.Context, name string, userID uint, deptID uint) ([]*dtos.Category, *common.ErrorCodeMessage) {
	if deptID == 0 {
		return nil, nil
	}
	name = strings.ToLower(name)
	cat, err := c.categoryRepo.FindCategoryByNameLike(name, deptID)
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
	return mapper.CategoriesModelToDTO(cat), nil
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
	isDeptManager, isOrgManager, _ := c.CheckUserRoleInOrganization(category.OrganizationID, userID, category.DepartmentID)
	if !isOrgManager && !isDeptManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if req.Name != nil {
		categoryByName, err := c.categoryRepo.FindCategoryByName(*req.Name, category.DepartmentID)
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

func (c *CategoryService) UpdateCategoryStatus(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryStatusRequest) *common.ErrorCodeMessage {
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
	isDeptManager, isOrgManager, _ := c.CheckUserRoleInOrganization(category.OrganizationID, userID, category.DepartmentID)
	if !isOrgManager && !isDeptManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if err := c.categoryRepo.UpdateCategoryStatusByID(catID, *req.Status); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	// update document status
	if err := c.documentRepo.UpdateDocumentStatusByCategoryID(catID, *req.Status); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (c *CategoryService) UpdateDepartmentStatus(ctx context.Context, userID uint, deptID uint, req *request.UpdateDepartmentStatusRequest) *common.ErrorCodeMessage {
	deparment, err := c.categoryRepo.FindDepartmentByID(deptID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if deparment == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCategoryNotFound,
			Message:     common.ErrMessageCategoryNotFound,
		}
	}
	isDeptManager, _, _ := c.CheckUserRoleInOrganization(deparment.OrganizationID, userID, deptID)
	if !isDeptManager {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}

	if err := c.categoryRepo.UpdateDepartmentStatusByID(deptID, *req.Status); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	if err := c.categoryRepo.UpdateCategoriesStatusByDepartmentID(deptID, *req.Status); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	if err := c.documentRepo.UpdateDocumentStatusByDepartmentID(deptID, *req.Status); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}

	return nil
}

func (c *CategoryService) CheckUserRoleInOrganization(orgID uint, userID uint, depID uint) (bool, bool, *common.ErrorCodeMessage) {
	isOrgManager := false
	isDeptManager := false
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return false, false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if user.OrganizationID == 0 || user.OrganizationID != orgID {
		return false, false, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeCannotAccessToOrganization,
			Message:     common.ErrMessageCannotAccessToOrganization,
		}
	}
	if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsOrganizationManager {
		isOrgManager = true
	}
	if user.OrganizationID != 0 && user.OrganizationID == orgID && user.IsDeptManager && user.DeptID == depID {
		isDeptManager = true
	}
	return isDeptManager, isOrgManager, nil
}

func (c *CategoryService) buildUpdateCategory(existedCategory *model2.Category, req *request.UpdateCategoryRequest) *model2.Category {
	var category model2.Category
	if req == nil {
		return existedCategory
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
