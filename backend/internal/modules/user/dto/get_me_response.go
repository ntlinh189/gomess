package dto

type GetMeResponse struct {
	ID      int64  `json:"id"`
	Account string `json:"account"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
}
