package dto

type SearchResponse struct {
	ID       int64  `json:"id"`
	Provider string `json:"provider"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}
