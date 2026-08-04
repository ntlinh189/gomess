package auth

import (
	"errors"
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
	cfg     config.ConfigInterface
}

func NewHandler(service ServiceInterface, cfg config.ConfigInterface) *Handler {
	return &Handler{service: service, cfg: cfg}
}

// Login godoc
//
//	@Summary		Login with OAuth
//	@Description	Login or register user using OAuth provider
//	@Tags			Authentication
//
//	@Accept			json
//	@Produce		json
//
//	@Param			provider	path		string				true	"OAuth provider"
//	@Param			request		body		dto.LoginRequest	true	"Login request"
//
//	@Success		200			{object}	dto.LoginResponse
//
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	provider := c.Param("provider")

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, refreshToken, err := h.service.Login(provider, req.Token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login failed"})
		return
	}

	h.setRefreshCookie(c, refreshToken)

	c.JSON(http.StatusOK, result)
}

// Refresh godoc
//
//	@Summary		Refresh access token
//	@Description	Generate a new access token using refresh token stored in cookie
//	@Tags			Authentication
//
//	@Produce		json
//
//	@Success		200	{object}	dto.RefreshResponse
//
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired cookie refresh token"})
		return
	}

	result, newRefreshToken, err := h.service.Refresh(refreshToken)

	if err != nil {
		h.clearRefreshCookie(c)

		if errors.Is(err, ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	h.setRefreshCookie(c, newRefreshToken)

	c.JSON(http.StatusOK, result)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout current user and invalidate refresh token
//	@Tags			Authentication
//
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Router			/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil {
		h.service.Logout(refreshToken)
	}

	h.clearRefreshCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) sameSiteMode() http.SameSite {
	if h.cfg.IsProduction() {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

func (h *Handler) setRefreshCookie(c *gin.Context, refreshToken string) {
	c.SetSameSite(h.sameSiteMode())
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/",
		"",
		h.cfg.IsProduction(),
		true,
	)
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(h.sameSiteMode())
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		h.cfg.IsProduction(),
		true,
	)
}