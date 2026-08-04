package models

import "time"

type FriendRequest struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
