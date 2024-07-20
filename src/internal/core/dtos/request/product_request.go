package request

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Type        *string  `json:"type"`
	Price       *float64 `json:"price"`
	Quantity    *int64   `json:"quantity"`
}

type ListProductRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
