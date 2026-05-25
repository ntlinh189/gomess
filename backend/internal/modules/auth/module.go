package auth

import (
	"gomess/internal/modules"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
}

func NewModule(ctx *modules.ModuleContext) *Module {

	handler := NewHandler(NewService(NewRepository(ctx.DB), ctx.JWT, ctx.Cfg, ctx.Redis), ctx.Cfg)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.POST("/:provider", m.handler.Login)
	auth.POST("/refresh", m.handler.Refresh)
	auth.POST("/logout", m.handler.Logout)
}