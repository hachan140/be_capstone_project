package dtos

type PrivateDocs struct {
	ID     uint `json:"id,omitempty"`
	DocID  uint `json:"doc_id,omitempty"`
	UserID uint `json:"user_id,omitempty"`
	Status int  `json:"status,omitempty"`
}
