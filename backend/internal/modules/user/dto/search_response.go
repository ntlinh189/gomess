package dto

type SearchResponse struct {
	ID       int64  `json:"id"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}
