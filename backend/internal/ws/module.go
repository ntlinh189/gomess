package ws

import (
	"gomess/internal/modules"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule(ctx *modules.ModuleContext, hub *Hub) *Module {
	return &Module{handler: NewHandler(hub, ctx.JWT, ctx.Cfg)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	ws := rg.Group("/ws")
	
	ws.GET("", m.handler.Serve)
}