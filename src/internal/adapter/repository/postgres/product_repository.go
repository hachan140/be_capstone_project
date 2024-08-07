package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
)

type IProductRepository interface {
	CreateProduct(product *model.Product) error
	GetProductByID(ID uint) (*model.Product, error)
	GetProductsByName(name string) ([]*model.Product, error)
	GetListProducts(limit int, offset int) ([]*model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id uint) error
}

type ProductRepository struct {
	storage *storage.Database
}

func NewProductRepository(storage *storage.Database) IProductRepository {
	return &ProductRepository{storage: storage}
}

func (p *ProductRepository) CreateProduct(product *model.Product) error {
	err := p.storage.Model(product).Create(product).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetProductByID(ID uint) (*model.Product, error) {
	var product *model.Product
	err := p.storage.Raw("select * from products where id = ?", ID).Scan(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductRepository) GetProductsByName(name string) ([]*model.Product, error) {
	var products []*model.Product
	err := p.storage.Raw("select * from products where LOWER(name) like ?", "%"+name+"%").Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductRepository) GetListProducts(limit int, offset int) ([]*model.Product, error) {
	var products []*model.Product
	err := p.storage.Raw("select * from products order by name asc  limit ? offset ?", limit, offset).Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductRepository) UpdateProduct(product *model.Product) error {
	err := p.storage.Model(product).Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	err := p.storage.Exec("delete from products where id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
