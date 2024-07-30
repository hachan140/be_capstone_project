package request

import "errors"

type SearchHistoryRequest struct {
	Keyword string `json:"keyword"`
}

type SaveSearchHistoryRequest struct {
	Keyword string `json:"keyword"`
	UserID  uint   `json:"user_id"`
	Type    int    `json:"type"`
}

func (s *SaveSearchHistoryRequest) Validate() error {
	if s.Keyword == "" {
		return errors.New("Keyword can't be null")
	}
	return nil
}
