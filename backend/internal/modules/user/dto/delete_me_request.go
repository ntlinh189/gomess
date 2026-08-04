package dto

type DeleteMeRequest struct {
	Confirm bool `json:"confirm" binding:"required"`
}