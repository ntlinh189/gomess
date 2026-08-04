package message

import (
	"gomess/internal/middleware"
	"gomess/internal/modules"
	"gomess/pkg/jwt"
	"gomess/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
	jwt     jwt.JWTInterface
}

func NewModule(
	ctx *modules.ModuleContext, 
	friendRepo FriendRepositoryInterface, 
	storage storage.StorageInterface,
	broadcaster Broadcaster,
) *Module {
	repo := NewRepository(ctx.DB)
	service := NewService(repo, friendRepo, storage, ctx.DB, broadcaster)
	handler := NewHandler(service)

	return &Module{handler: handler, jwt: ctx.JWT}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	messages := rg.Group("/messages")

	protected := messages.Group("")
	protected.Use(middleware.Auth(m.jwt))

	protected.POST("", m.handler.SendMessage)
	protected.GET("/:friendId", m.handler.GetHistory)
	protected.DELETE("/:id", m.handler.DeleteForMe)
	protected.POST("/:id/revoke", m.handler.RevokeMessage)
}
