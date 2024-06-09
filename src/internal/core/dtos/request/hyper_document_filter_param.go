package request

import "time"

type HyperDocumentFilterParam struct {
	Page            int       `form:"page"`
	PageSize        int       `form:"page_size"`
	Title           string    `form:"document_name"`
	Type            string    `form:"type"`
	CreatedFromDate time.Time `form:"created_from_date"`
	CreatedToDate   time.Time `form:"created_to_date"`
	CreatedBy       string    `form:"created_by"`
}
