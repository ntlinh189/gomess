package dto

import "time"

type AttachmentResponse struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type MessageResponse struct {
	ID          int64                `json:"id"`
	SenderID    int64                `json:"sender_id"`
	ReceiverID  int64                `json:"receiver_id"`
	Content     string               `json:"content"`
	Attachments []AttachmentResponse `json:"attachments"`
	CreatedAt   time.Time            `json:"created_at"`
	Revoked     bool                 `json:"revoked"`
}
