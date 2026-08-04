package dto

type FriendResponse struct {
	ID       int64  `json:"id"`
	Provider string `json:"provider"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}