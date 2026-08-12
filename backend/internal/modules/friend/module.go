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

	friends.Use(middleware.Auth(m.jwt))

	friends.GET("", m.handler.GetFriends)
	friends.DELETE("/:id", m.handler.DeleteFriend)

	friends.POST("/requests", m.handler.SendRequest)
	friends.GET("/requests/received", m.handler.GetReceivedRequests)
	friends.POST("/requests/:id/accept", m.handler.AcceptRequest)
	friends.POST("/requests/:id/reject", m.handler.RejectRequest)
}
