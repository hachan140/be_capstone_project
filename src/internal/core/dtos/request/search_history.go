package request

type SearchHistoryRequest struct {
	Keyword string `json:"keyword"`
}

type SaveSearchHistoryRequest struct {
	Keyword string `json:"keyword"`
	UserID  uint   `json:"user_id"`
	Type    int    `json:"type"`
}
