package ws

import (
	"gomess/internal/config"
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/pkg/jwt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub      *Hub
	jwt      jwt.JWTInterface
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, jwt jwt.JWTInterface, cfg config.ConfigInterface) *Handler {
	return &Handler{
		hub: hub,
		jwt: jwt,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get(context.OriginKey)
				if origin == "" {
					return true
				}
				return slices.Contains(cfg.GetClientOrigins(), origin)
			},
		},
	}
}

func (h *Handler) Serve(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	userID, err := h.jwt.ParseAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.FromGin(c).Error("websocket upgrade failed", "error", err)
		return
	}

	client := newClient(h.hub, conn, userID)
	h.hub.register(client)

	go client.writePump()
	go client.readPump()
}
