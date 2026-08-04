package upload

import (
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/internal/modules/upload/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	Presign(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// Presign godoc
//
//	@Summary		Get presigned upload URL
//	@Description	Get a short-lived URL to upload a file directly to object storage
//	@Tags			Upload
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.PresignRequest	true	"File info"
//	@Success		200		{object}	dto.PresignResponse
//	@Router			/uploads/presign [post]
func (h *Handler) Presign(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	var req dto.PresignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.PresignUpload(userID, req.FileName)
	if err != nil {
		logger.FromGin(c).Error("presign upload error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}