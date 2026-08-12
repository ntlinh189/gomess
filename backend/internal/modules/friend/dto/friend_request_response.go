package dto

import "time"

type FriendRequestResponse struct {
	ID         int64          `json:"id"`
	ReceiverID int64          `json:"receiver_id"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	Sender     FriendResponse `json:"sender"`
}
