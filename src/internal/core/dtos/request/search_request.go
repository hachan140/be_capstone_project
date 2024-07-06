package request

type SearchAndOrNotRequest struct {
	AndKeyWords []string `json:"and_keywords"`
	OrKeyWords  []string `json:"or_keywords"`
	NotKeywords []string `json:"not_keywords"`
}
