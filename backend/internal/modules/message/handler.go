package message

import (
	"errors"
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/internal/modules/message/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	SendMessage(c *gin.Context)
	GetHistory(c *gin.Context)
	DeleteForMe(c *gin.Context)
	RevokeMessage(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// SendMessage godoc
//
//	@Summary		Send a message
//	@Description	Send a message to a friend (only friends can message each other)
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.SendMessageRequest	true	"Message"
//	@Success		201		{object}	dto.MessageResponse
//	@Router			/messages [post]
func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	message, err := h.service.SendMessage(userID, req.ReceiverID, req.Content, req.Attachments)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetHistory godoc
//
//	@Summary		Get message history with a friend
//	@Tags			Message
//	@Produce		json
//	@Security		BearerAuth
//	@Param			friendId	path	int	true	"Friend user ID"
//	@Param			before_id	query	int	false	"Get messages older than this ID"
//	@Param			limit		query	int	false	"Limit (default 20, max 100)"
//	@Success		200			{array}	dto.MessageResponse
//	@Router			/messages/{friendId} [get]
func (h *Handler) GetHistory(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	friendID, err := strconv.ParseInt(c.Param("friendId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friend id"})
		return
	}

	beforeID, _ := strconv.ParseInt(c.Query("before_id"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit"))

	messages, err := h.service.GetHistory(userID, friendID, beforeID, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// DeleteForMe godoc
//
//	@Summary		Delete a message for myself only
//	@Tags			Message
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Message ID"
//	@Success		200
//	@Router			/messages/{id} [delete]
func (h *Handler) DeleteForMe(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	if err := h.service.DeleteForMe(userID, messageID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}

// RevokeMessage godoc
//
//	@Summary		Revoke a message for everyone
//	@Tags			Message
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Message ID"
//	@Success		200
//	@Router			/messages/{id}/revoke [post]
func (h *Handler) RevokeMessage(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	if err := h.service.RevokeMessage(userID, messageID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message revoked"})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFriend):
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only message friends"})
	case errors.Is(err, ErrCannotMessageYourself):
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot send message to yourself"})
	case errors.Is(err, ErrEmptyMessage):
		c.JSON(http.StatusBadRequest, gin.H{"error": "message must have content or at least one attachment"})
	case errors.Is(err, ErrTooManyAttachments):
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many attachments"})
	case errors.Is(err, ErrAttachmentNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more attachments were not uploaded"})
	case errors.Is(err, ErrMessageNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, ErrAlreadyRevoked):
		c.JSON(http.StatusConflict, gin.H{"error": "message already revoked"})
	case errors.Is(err, ErrRevokeWindowExpired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "revoke window expired"})
	default:
		logger.FromGin(c).Error("message module error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
