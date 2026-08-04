package ws

import (
	"gomess/internal/modules"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule(ctx *modules.ModuleContext, hub *Hub) *Module {
	return &Module{handler: NewHandler(hub, ctx.JWT)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ws", m.handler.Serve)
}