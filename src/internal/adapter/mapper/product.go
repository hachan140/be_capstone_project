package mapper

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/dtos"
)

func ProductModelToDTO(product *model.Product) *dtos.Product {
	if product == nil {
		return nil
	}
	return &dtos.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Type:        product.Type,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ProductModelsToDTOs(products []*model.Product) []*dtos.Product {
	if products == nil {
		return nil
	}
	productDTOs := make([]*dtos.Product, 0)
	for _, p := range products {
		productDTO := &dtos.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Type:        p.Type,
			Price:       p.Price,
			Quantity:    p.Quantity,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
		productDTOs = append(productDTOs, productDTO)
	}
	return productDTOs
}
