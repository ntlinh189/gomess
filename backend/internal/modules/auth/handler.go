package auth

import (
	"gomess/internal/modules/auth/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	Login(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	provider := c.Param("provider")

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Login(provider, req.Token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}