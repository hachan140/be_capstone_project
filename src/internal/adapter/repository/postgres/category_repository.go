package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type ICategoryRepository interface {
	CreateCategory(category *model.Category) error
	UpdateCategory(category *model.Category) error
	ListCategoryByOrganization(orgID uint, limit int, offset int) ([]*model.Category, error)
	FindCategoryByID(catID uint) (*model.Category, error)
	FindCategoryByName(name string) (*model.Category, error)
}

type CategoryRepository struct {
	storage *storage.Database
}

func NewCategoryRepository(storage *storage.Database) ICategoryRepository {
	return &CategoryRepository{storage: storage}
}

func (c *CategoryRepository) CreateCategory(category *model.Category) error {
	err := c.storage.Model(category).Create(&category).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateCategory(category *model.Category) error {
	err := c.storage.Model(category).Where("id = ?", category.ID).Updates(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) ListCategoryByOrganization(orgID uint, limit int, offset int) ([]*model.Category, error) {
	var categories []*model.Category
	err := c.storage.Raw("select * from categories where organization_id = ? order by created_at desc limit ? offset ?", orgID, limit, offset).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *CategoryRepository) FindCategoryByID(catID uint) (*model.Category, error) {
	var category *model.Category
	err := c.storage.Raw("select * from categories where id = ?", catID).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoryRepository) FindCategoryByName(name string) (*model.Category, error) {
	var category *model.Category
	err := c.storage.Raw("select * from categories where name = ?", name).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}
