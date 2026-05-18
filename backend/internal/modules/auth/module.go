package auth

import (
	"gomess/internal/modules"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
}

func NewModule(ctx *modules.ModuleContext) *Module {

	handler := NewHandler(NewService(NewRepository(ctx.DB), ctx.JWT, ctx.Cfg, ctx.Redis))

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.POST("/:provider", m.handler.Login)
}