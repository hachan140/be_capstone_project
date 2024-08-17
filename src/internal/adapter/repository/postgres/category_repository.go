package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type ICategoryRepository interface {
	CreateCategory(category *model.Category) error
	UpdateCategory(category *model.Category) error
	UpdateDepartment(department *model.Department) error
	ListCategoryByDepartment(orgID uint, limit int, offset int) ([]*model.Category, error)
	ListCategoryByDepartmentID(orgID uint) ([]*model.Category, error)
	FindCategoryByID(catID uint) (*model.Category, error)
	FindDepartmentByID(deptID uint) (*model.Department, error)
	FindCategoryByName(name string, deptID uint) (*model.Category, error)
	FindCategoryByNameLike(name string, deptID uint) ([]*model.Category, error)
	UpdateCategoriesStatusByDepartmentID(deptID uint, status int) error
	UpdateCategoriesStatusByOrganizationID(orgID uint, status int) error
	UpdateDepartmentsStatusByOrganizationID(orgID uint, status int) error
	UpdateDepartmentStatusByID(deptID uint, status int) error
	UpdateCategoryStatusByID(catID uint, status int) error
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

func (c *CategoryRepository) UpdateDepartment(department *model.Department) error {
	err := c.storage.Model(department).Where("id = ?", department.ID).Updates(department).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateDepartmentStatusByID(deptID uint, status int) error {
	err := c.storage.Exec("update departments set status = ? where id = ?", status, deptID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateCategoriesStatusByOrganizationID(orgID uint, status int) error {
	err := c.storage.Exec("update categories set status = ? where organization_id = ?", status, orgID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateCategoriesStatusByDepartmentID(deptID uint, status int) error {
	err := c.storage.Exec("update categories set status = ? where department_id = ?", status, deptID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateCategoryStatusByID(catID uint, status int) error {
	err := c.storage.Exec("update categories set status = ? where id = ?", status, catID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) UpdateDepartmentsStatusByOrganizationID(orgID uint, status int) error {
	err := c.storage.Exec("update departments set status = ? where organization_id = ?", status, orgID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) ListCategoryByDepartment(orgID uint, limit int, offset int) ([]*model.Category, error) {
	var categories []*model.Category
	err := c.storage.Raw("select * from categories where department_id = ? and status = 1 order by created_at desc limit ? offset ?", orgID, limit, offset).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *CategoryRepository) ListCategoryByDepartmentID(orgID uint) ([]*model.Category, error) {
	var categories []*model.Category
	err := c.storage.Raw("select * from categories where department_id = ?", orgID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *CategoryRepository) FindCategoryByID(catID uint) (*model.Category, error) {
	var category *model.Category
	err := c.storage.Raw("select * from categories where id = ? and status = 1", catID).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoryRepository) FindDepartmentByID(deptID uint) (*model.Department, error) {
	var department *model.Department
	err := c.storage.Raw("select * from departments where id = ?", deptID).Scan(&department).Error
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (c *CategoryRepository) FindCategoryByName(name string, deptID uint) (*model.Category, error) {
	var category *model.Category
	err := c.storage.Raw("select * from categories where name = ? and department_id = ? and status = 1", name, deptID).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoryRepository) FindCategoryByNameLike(name string, deptID uint) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	err := c.storage.Raw("select * from categories where lower(name) like ? and department_id = ?", "%"+name+"%", deptID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
