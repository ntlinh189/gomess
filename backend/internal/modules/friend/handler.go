package friend

import (
	"errors"
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/internal/modules/friend/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	SendRequest(c *gin.Context)
	AcceptRequest(c *gin.Context)
	RejectRequest(c *gin.Context)
	DeleteFriend(c *gin.Context)
	GetFriends(c *gin.Context)
	GetReceivedRequests(c *gin.Context)
	GetSentRequests(c *gin.Context)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// SendRequest godoc
//
//	@Summary		Send friend request
//	@Description	Send a friend request to another user (auto-accepts if that user already sent one to you)
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	dto.SendRequestRequest	true	"Receiver"
//	@Success		201
//	@Router			/friends/requests [post]
func (h *Handler) SendRequest(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	var req dto.SendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.SendRequest(userID, req.ReceiverID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "friend request sent"})
}

// AcceptRequest godoc
//
//	@Summary		Accept friend request
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Friend request ID"
//	@Success		200
//	@Router			/friends/requests/{id}/accept [post]
func (h *Handler) AcceptRequest(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	if err := h.service.AcceptRequest(userID, requestID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

// RejectRequest godoc
//
//	@Summary		Reject friend request
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Friend request ID"
//	@Success		200
//	@Router			/friends/requests/{id}/reject [post]
func (h *Handler) RejectRequest(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	if err := h.service.RejectRequest(userID, requestID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request rejected"})
}

// DeleteFriend godoc
//
//	@Summary		Unfriend
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Friend user ID"
//	@Success		200
//	@Router			/friends/{id} [delete]
func (h *Handler) DeleteFriend(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friend id"})
		return
	}

	if err := h.service.DeleteFriend(userID, friendID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend removed"})
}

// GetFriends godoc
//
//	@Summary		Get friend list
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	dto.FriendResponse
//	@Router			/friends [get]
func (h *Handler) GetFriends(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	friends, err := h.service.GetFriends(userID)
	if err != nil {
		logger.FromGin(c).Error("get friends error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]dto.FriendResponse, 0, len(friends))
	for _, friend := range friends {
		resp = append(resp, dto.FriendResponse{
			ID:       friend.ID,
			Provider: friend.Provider,
			Account:  friend.Account,
			Name:     friend.Name,
			Avatar:   friend.Avatar,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetReceivedRequests godoc
//
//	@Summary		Get received friend requests
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	dto.FriendRequestResponse
//	@Router			/friends/requests/received [get]
func (h *Handler) GetReceivedRequests(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	requests, err := h.service.GetReceivedRequests(userID)
	if err != nil {
		logger.FromGin(c).Error("get received requests error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]dto.FriendRequestResponse, 0, len(requests))
	for _, request := range requests {
		resp = append(resp, dto.FriendRequestResponse{
			ID:         request.ID,
			SenderID:   request.SenderID,
			ReceiverID: request.ReceiverID,
			Status:     request.Status,
			CreatedAt:  request.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetSentRequests godoc
//
//	@Summary		Get sent friend requests
//	@Tags			Friend
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	dto.FriendRequestResponse
//	@Router			/friends/requests/sent [get]
func (h *Handler) GetSentRequests(c *gin.Context) {
	userID := c.GetInt64(context.UserIDKey)

	requests, err := h.service.GetSentRequests(userID)
	if err != nil {
		logger.FromGin(c).Error("get sent requests error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]dto.FriendRequestResponse, 0, len(requests))
	for _, request := range requests {
		resp = append(resp, dto.FriendRequestResponse{
			ID:         request.ID,
			SenderID:   request.SenderID,
			ReceiverID: request.ReceiverID,
			Status:     request.Status,
			CreatedAt:  request.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrCannotAddYourself):
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add yourself"})
	case errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, ErrAlreadyFriend):
		c.JSON(http.StatusConflict, gin.H{"error": "already friend"})
	case errors.Is(err, ErrRequestExists):
		c.JSON(http.StatusConflict, gin.H{"error": "friend request already exists"})
	case errors.Is(err, ErrRequestNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, ErrNotFriend):
		c.JSON(http.StatusNotFound, gin.H{"error": "not friend"})
	default:
		logger.FromGin(c).Error("friend module error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}