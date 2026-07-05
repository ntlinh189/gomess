package modules

import (
	"gomess/internal/config"
	"gomess/internal/database"
	"gomess/internal/redis"
	"gomess/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type ModuleInterface interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

type ModuleContext struct {
	DB    database.DatabaseInterface
	JWT   jwt.JWTInterface
	Cfg   config.ConfigInterface
	Redis redis.RedisInterface
}
