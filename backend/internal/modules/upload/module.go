package upload

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

func NewModule(ctx *modules.ModuleContext, storage storage.StorageInterface) *Module {
	service := NewService(storage)
	handler := NewHandler(service)

	return &Module{handler: handler, jwt: ctx.JWT}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	uploads := rg.Group("/uploads")
	uploads.Use(middleware.Auth(m.jwt))

	uploads.POST("/presign", m.handler.Presign)
}