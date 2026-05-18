package user

import (
	"gomess/internal/modules"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
}

func NewModule(ctx *modules.ModuleContext) *Module {

	handler := NewHandler(NewService(NewRepository(ctx.DB)))

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.GET("ping", func(c *gin.Context) { c.JSON(http.StatusOK, "pong") })
}
