package request

type SearchAndOrNotRequest struct {
	AndKeyWords []string `json:"and_key_words"`
	OrKeyWords  []string `json:"or_key_words"`
	NotKeywords []string `json:"not_keywords"`
}
