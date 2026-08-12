package models

import "time"

type FriendRequest struct {
	ID         int64
	Sender     User
	ReceiverID int64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
