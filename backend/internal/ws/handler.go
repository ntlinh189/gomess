package ws

import (
	"gomess/internal/context"
	"gomess/internal/logger"
	"gomess/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool {
		origin := r.Header.Get(context.OriginKey)
		if origin == "" {
			return true
		}
		/// TODO: Change to client url
		return origin == "http://localhost:3000"
	},
}

type Handler struct {
	hub *Hub
	jwt jwt.JWTInterface
}

func NewHandler(hub *Hub, jwt jwt.JWTInterface) *Handler {
	return &Handler{
		hub: hub,
		jwt: jwt,
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.FromGin(c).Error("websocket upgrade failed", "error", err)
		return
	}

	client := newClient(h.hub, conn, userID)
	h.hub.register(client)

	go client.writePump()
	go client.readPump()
}
