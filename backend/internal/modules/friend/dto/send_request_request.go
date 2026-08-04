package dto

type SendRequestRequest struct {
	ReceiverID int64 `json:"receiver_id" binding:"required"`
}