package user

import (
	"gomess/internal/modules/user/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	GetMe(c *gin.Context)
	Search(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	users, err := h.service.Search(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
