package request

type GetListCategoryRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
