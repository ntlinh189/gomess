package dto

type AttachmentInput struct {
	ObjectKey string `json:"object_key" binding:"required"`
	FileName  string `json:"file_name" binding:"required"`
}

type SendMessageRequest struct {
	ReceiverID  int64             `json:"receiver_id" binding:"required"`
	Content     string            `json:"content" binding:"max=2000"`
	Attachments []AttachmentInput `json:"attachments" binding:"max=10,dive"`
}
