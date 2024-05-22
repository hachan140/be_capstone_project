package services

import (
	"be-capstone-project/src/internal/adapter/mapper"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	model2 "be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"be-capstone-project/src/internal/core/dtos/request"
	"context"
	"errors"
	"time"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, userID uint, req *request.CreateCategoryRequest) error
	ListCategories(ctx context.Context, orgID uint, userID uint) ([]*dtos.Category, error)
	GetCategoryByID(ctx context.Context, id uint, userID uint) (*dtos.Category, error)
	UpdateCategoryByID(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryRequest) error
}

type CategoryService struct {
	categoryRepo   postgres.ICategoryRepository
	userRepository postgres.IUserRepository
}

func NewCategoryService(categoryRepo postgres.ICategoryRepository, userRepo postgres.IUserRepository) ICategoryService {
	return &CategoryService{categoryRepo: categoryRepo, userRepository: userRepo}
}

func (c *CategoryService) CreateCategory(ctx context.Context, userID uint, req *request.CreateCategoryRequest) error {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return err
	}
	isOrgManager, err := c.CheckUserRoleInOrganization(req.OrganizationID, user.ID)
	if err != nil {
		return err
	}
	if !isOrgManager || err != nil {
		return errors.New(common.ErrMessageCannotAccessToOrganization)
	}
	catExisted, err := c.categoryRepo.FindCategoryByName(req.Name)
	if err != nil {
		return err
	}
	if catExisted != nil {
		return errors.New(common.ErrMessageCategoryExisted)
	}
	model := &model2.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryID: req.ParentID,
		OrganizationID:   req.OrganizationID,
		Status:           1,
		CreatedBy:        req.CreatedBy,
		CreatedAt:        time.Now(),
	}
	if err := c.categoryRepo.CreateCategory(model); err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) ListCategories(ctx context.Context, orgID uint, userID uint) ([]*dtos.Category, error) {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return nil, err
	}
	if user.OrganizationID != 0 && user.OrganizationID != orgID {
		return nil, errors.New(common.ErrMessageUserAlreadyInOtherOrganization)
	}
	categories, err := c.categoryRepo.ListCategoryByOrganization(orgID)
	if err != nil {
		return nil, err
	}
	var categoryRes []*dtos.Category
	for _, c := range categories {
		cRes := &dtos.Category{
			ID:               c.ID,
			Name:             c.Name,
			Description:      c.CreatedBy,
			ParentCategoryID: c.ParentCategoryID,
			OrganizationID:   c.OrganizationID,
			Status:           c.Status,
			CreatedBy:        c.CreatedBy,
			CreatedAt:        time.Now(),
			UpdatedAt:        c.UpdatedAt,
		}
		categoryRes = append(categoryRes, cRes)
	}
	return categoryRes, nil
}

func (c *CategoryService) GetCategoryByID(ctx context.Context, id uint, userID uint) (*dtos.Category, error) {

	cat, err := c.categoryRepo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New(common.ErrMessageCategoryNotFound)
	}
	_, err = c.CheckUserRoleInOrganization(cat.OrganizationID, userID)
	if err != nil {
		return nil, err
	}
	return mapper.CategoryModelToDTO(cat), nil
}

func (c *CategoryService) UpdateCategoryByID(ctx context.Context, userID uint, catID uint, req *request.UpdateCategoryRequest) error {
	category, err := c.categoryRepo.FindCategoryByID(catID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New(common.ErrMessageCategoryNotFound)
	}
	isOrgManager, err := c.CheckUserRoleInOrganization(category.OrganizationID, userID)
	if !isOrgManager || err != nil {
		return errors.New(common.ErrMessageCannotAccessToOrganization)
	}
	if req.Name != nil {
		categoryByName, err := c.categoryRepo.FindCategoryByName(*req.Name)
		if err != nil {
			return err
		}
		if categoryByName != nil && categoryByName.Name != category.Name {
			return errors.New(common.ErrMessageCategoryExisted)
		}
	}

	categoryToUpdate := c.buildUpdateCategory(category, req)
	if err := c.categoryRepo.UpdateCategory(categoryToUpdate); err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) CheckUserRoleInOrganization(orgID uint, userID uint) (bool, error) {
	user, err := c.userRepository.FinduserByID(userID)
	if err != nil {
		return false, err
	}
	if user.OrganizationID == 0 || user.OrganizationID != orgID {
		return false, errors.New(common.ErrMessageCannotAccessToOrganization)
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
