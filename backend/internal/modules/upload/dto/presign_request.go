package dto

type PresignRequest struct {
	FileName string `json:"file_name" binding:"required"`
}