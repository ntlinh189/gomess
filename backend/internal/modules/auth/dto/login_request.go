package dto

type LoginRequest struct {
	Token string `json:"token" binding:"required"`
}