package user

import (
	"gomess/internal/middleware"
	"gomess/internal/modules"
	"gomess/internal/redis"
	"gomess/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler HandlerInterface
	jwt     jwt.JWTInterface
	redis   redis.RedisInterface
}

func NewModule(ctx *modules.ModuleContext, repo RepositoryInterface) *Module {
	handler := NewHandler(NewService(repo), ctx.Cfg)

	return &Module{handler: handler, jwt: ctx.JWT, redis: ctx.Redis}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	protected := user.Group("/")
	protected.Use(middleware.Auth(m.jwt))

	protected.GET("/me", m.handler.GetMe)
	protected.GET(
		"/search", 
		middleware.RateLimit(m.redis, "search", 30, time.Minute, middleware.RateLimitByUser),
		m.handler.Search,
	)
	protected.DELETE("/me", m.handler.DeleteMe)
}
