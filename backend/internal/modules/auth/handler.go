package auth

import (
	"gomess/internal/config"
	"gomess/internal/modules/auth/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	Logout(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
	cfg config.ConfigInterface
}

func NewHandler(service ServiceInterface, cfg config.ConfigInterface) *Handler {
	return &Handler{service: service, cfg: cfg,}
}

func (h *Handler) Login(c *gin.Context) {
	provider := c.Param("provider")

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, refreshToken, err := h.service.Login(provider, req.Token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(func() http.SameSite {
		if h.cfg.IsProduction() {
			return http.SameSiteNoneMode
		}
		return http.SameSiteLaxMode
	}())

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((30*24*time.Hour).Seconds()),
		"/",
		"",
		h.cfg.IsProduction(),
		true,
	)

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Refresh(refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil {
		h.service.Logout(refreshToken)
	}

	c.SetSameSite(func() http.SameSite {
		if h.cfg.IsProduction() {
			return http.SameSiteNoneMode
		}
		return http.SameSiteLaxMode
	}())

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		h.cfg.IsProduction(),
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out",})
}