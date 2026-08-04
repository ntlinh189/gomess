package auth

import (
	"gomess/internal/middleware"
	"gomess/internal/modules"
	"gomess/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
	redis   redis.RedisInterface
}

func NewModule(ctx *modules.ModuleContext) *Module {

	handler := NewHandler(NewService(NewRepository(ctx.DB), ctx.JWT, ctx.Cfg, ctx.Redis), ctx.Cfg)

	return &Module{handler: handler, redis: ctx.Redis}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.POST(
		"/:provider", 
		middleware.RateLimit(m.redis, "login", 5, time.Minute, middleware.RateLimitByIP),
		m.handler.Login,
	)
	auth.POST("/refresh", m.handler.Refresh)
	auth.POST("/logout", m.handler.Logout)
}
