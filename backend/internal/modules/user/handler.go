package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	GetMe(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetInt64("userID")

	user, err := h.service.GetMe(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
