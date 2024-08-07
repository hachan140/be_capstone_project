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
	"time"
)

type IProductService interface {
	CreateProduct(ctx context.Context, userID uint, req *request.CreateProductRequest) *common.ErrorCodeMessage
	GetProductByID(ctx context.Context, ID uint) (*dtos.Product, *common.ErrorCodeMessage)
	GetProductsByName(ctx context.Context, name string) ([]*dtos.Product, *common.ErrorCodeMessage)
	GetListProducts(ctx context.Context, page int, pageSize int) ([]*dtos.Product, *common.ErrorCodeMessage)
	UpdateProduct(ctx context.Context, userID uint, productID uint, req *request.UpdateProductRequest) *common.ErrorCodeMessage
	DeleteProduct(ctx context.Context, userID uint, productID uint) *common.ErrorCodeMessage
}

type ProductService struct {
	productRepository postgres.IProductRepository
	userRepository    postgres.IUserRepository
}

func NewProductService(productRepo postgres.IProductRepository, userRepo postgres.IUserRepository) IProductService {
	return &ProductService{
		productRepository: productRepo,
		userRepository:    userRepo,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, userID uint, req *request.CreateProductRequest) *common.ErrorCodeMessage {
	if !p.isAdmin(userID) {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.DefaultMaxHeaderBytes,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	productToCreate := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Price:       req.Price,
		Quantity:    req.Quantity,
		CreatedAt:   time.Now(),
	}
	if err := p.productRepository.CreateProduct(productToCreate); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (p *ProductService) GetProductByID(ctx context.Context, ID uint) (*dtos.Product, *common.ErrorCodeMessage) {
	product, err := p.productRepository.GetProductByID(ID)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return mapper.ProductModelToDTO(product), nil
}

func (p *ProductService) GetProductsByName(ctx context.Context, name string) ([]*dtos.Product, *common.ErrorCodeMessage) {
	products, err := p.productRepository.GetProductsByName(strings.ToLower(name))
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return mapper.ProductModelsToDTOs(products), nil
}

func (p *ProductService) GetListProducts(ctx context.Context, page int, pageSize int) ([]*dtos.Product, *common.ErrorCodeMessage) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	products, err := p.productRepository.GetListProducts(pageSize, offset)
	if err != nil {
		return nil, &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return mapper.ProductModelsToDTOs(products), nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, userID uint, productID uint, req *request.UpdateProductRequest) *common.ErrorCodeMessage {
	if !p.isAdmin(userID) {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.DefaultMaxHeaderBytes,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	productExisted, err := p.productRepository.GetProductByID(productID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if productExisted == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeProductNotExist,
			Message:     common.ErrMessageProductNotExist,
		}
	}
	product := p.buildUpdateProductQuery(productExisted, req)
	if err := p.productRepository.UpdateProduct(product); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, userID uint, productID uint) *common.ErrorCodeMessage {
	if !p.isAdmin(userID) {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.DefaultMaxHeaderBytes,
			ServiceCode: common.ErrCodeUserDoesNotHavePermission,
			Message:     common.ErrMessageUserDoesNotHavePermission,
		}
	}
	productExisted, err := p.productRepository.GetProductByID(productID)
	if err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	if productExisted == nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: common.ErrCodeProductNotExist,
			Message:     common.ErrMessageProductNotExist,
		}
	}
	if err := p.productRepository.DeleteProduct(productID); err != nil {
		return &common.ErrorCodeMessage{
			HTTPCode:    http.StatusInternalServerError,
			ServiceCode: common.ErrCodeInternalError,
			Message:     err.Error(),
		}
	}
	return nil
}

func (p *ProductService) buildUpdateProductQuery(existed *model.Product, req *request.UpdateProductRequest) *model.Product {
	productToUpdate := model.Product{}
	if req == nil {
		return existed
	}
	if req.Name != nil {
		productToUpdate.Name = *req.Name
	} else {
		productToUpdate.Name = existed.Name
	}
	if req.Description != nil {
		productToUpdate.Description = *req.Description
	} else {
		productToUpdate.Description = existed.Description
	}
	if req.Quantity != nil {
		productToUpdate.Quantity = *req.Quantity
	} else {
		productToUpdate.Quantity = existed.Quantity
	}
	productToUpdate.ID = existed.ID
	productToUpdate.UpdatedAt = time.Now()
	return &productToUpdate
}

func (p *ProductService) isAdmin(userID uint) bool {
	user, err := p.userRepository.FinduserByID(userID)
	if err != nil {
		return false
	}
	if user.IsAdmin {
		return true
	}
	return false
}
