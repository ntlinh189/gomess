package user

import (
	"gomess/internal/modules"
	"gomess/internal/modules/middleware"
	"gomess/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
	jwt jwt.JWTInterface
}

func NewModule(ctx *modules.ModuleContext) *Module {

	handler := NewHandler(NewService(NewRepository(ctx.DB)))
	jwt := ctx.JWT

	return &Module{handler: handler, jwt: jwt}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.GET("ping", func(c *gin.Context) { c.JSON(http.StatusOK, "pong") })
	user.GET("/me", middleware.Auth(m.jwt), m.handler.GetMe)
}
