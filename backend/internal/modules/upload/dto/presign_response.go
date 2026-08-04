package dto

type PresignResponse struct {
	ObjectKey string `json:"object_key"`
	UploadURL string `json:"upload_url"`
}