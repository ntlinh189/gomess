package user

import (
	"errors"
	"gomess/internal/config"
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/internal/modules/user/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	GetMe(c *gin.Context)
	Search(c *gin.Context)
	DeleteMe(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
	cfg     config.ConfigInterface
}

func NewHandler(service ServiceInterface, cfg config.ConfigInterface) *Handler {
	return &Handler{service: service, cfg: cfg}
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Get information of authenticated user
//	@Tags			User
//
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Success		200	{object}	dto.GetMeResponse
//
//	@Router			/user/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	user, err := h.service.GetMe(userID)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		logger.FromGin(c).Error("get me error", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Search godoc
//
//	@Summary		Search users
//	@Description	Search users by provider and keyword
//	@Tags			User
//
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			provider	query	string	true	"Provider"
//	@Param			keyword		query	string	true	"Search keyword"
//	@Param			skip		query	int		false	"Skip"
//	@Param			limit		query	int		false	"Limit"
//
//	@Success		200			{array}	dto.SearchResponse
//
//	@Router			/user/search [get]
func (h *Handler) Search(c *gin.Context) {

	var req dto.SearchRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	users, err := h.service.Search(c.GetInt64(context.UserIDKey), &req)

	if err != nil {
		logger.FromGin(c).Error("search users error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteMe godoc
//
//	@Summary		Delete my account
//	@Description	Permanently delete the authenticated user's account and all related data (friends, messages, attachment metadata)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	dto.DeleteMeRequest	true	"Confirmation"
//	@Success		200
//	@Router			/user/me [delete]
func (h *Handler) DeleteMe(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	var req dto.DeleteMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "must confirm account deletion"})
		return
	}

	if err := h.service.DeleteMe(userID); err != nil {
		logger.FromGin(c).Error("delete me error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})

	// Should call /auth/logout in client side to clear refresh token
}
