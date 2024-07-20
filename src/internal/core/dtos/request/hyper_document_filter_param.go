package request

import "time"

type HyperDocumentFilterParam struct {
	Page            int       `form:"page"`
	PageSize        int       `form:"page_size"`
	Title           string    `form:"title"`
	Type            string    `form:"type"`
	CreatedFromDate time.Time `form:"created_from_date"  time_format:"2006-01-02"`
	CreatedToDate   time.Time `form:"created_to_date"  time_format:"2006-01-02"`
	CreatedBy       string    `form:"created_by"`
}
