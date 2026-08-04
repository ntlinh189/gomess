package friend

import (
	"gomess/internal/middleware"
	"gomess/internal/modules"
	"gomess/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
	jwt     jwt.JWTInterface
}

func NewModule(
	ctx *modules.ModuleContext, 
	repo RepositoryInterface,
	userRepo UserRepositoryInterface, 
	messageRepo MessageRepositoryInterface,
) *Module {
	service := NewService(repo, userRepo, messageRepo, ctx.DB)
	handler := NewHandler(service)

	return &Module{handler: handler, jwt: ctx.JWT}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	friends := rg.Group("/friends")

	protected := friends.Group("/")
	protected.Use(middleware.Auth(m.jwt))

	protected.GET("", m.handler.GetFriends)
	protected.DELETE("/:id", m.handler.DeleteFriend)

	protected.POST("/requests", m.handler.SendRequest)
	protected.GET("/requests/received", m.handler.GetReceivedRequests)
	protected.GET("/requests/sent", m.handler.GetSentRequests)
	protected.POST("/requests/:id/accept", m.handler.AcceptRequest)
	protected.POST("/requests/:id/reject", m.handler.RejectRequest)
}