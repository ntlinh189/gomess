package models

import "time"

type Message struct {
	ID          int64
	SenderID    int64
	ReceiverID  int64
	Content     string
	CreatedAt   time.Time
	RevokedAt   *time.Time
	Attachments []MessageAttachment
}

type MessageAttachment struct {
	ID        int64
	MessageID int64
	Type      string
	ObjectKey string
	FileName  string
	MimeType  string
	SizeBytes int64
	CreatedAt time.Time
}
