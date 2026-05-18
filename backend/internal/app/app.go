package app

import (
	"context"
	"gomess/internal/config"
	"gomess/internal/database"
	"gomess/internal/modules"
	"gomess/internal/modules/auth"
	"gomess/internal/modules/user"
	"gomess/internal/redis"
	"gomess/pkg/jwt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	router *gin.Engine
	config config.ConfigInterface
}

func NewApplication(cfg config.ConfigInterface) (*Application, error) {
	r := gin.Default()

	db, err := database.NewMySql(cfg)
	if err != nil {
		return nil, err
	}

	jwt := jwt.NewJWT(cfg.GetJWTSecret())

	redis := redis.NewRedis(cfg.GetRedisAddr())

	ctx := &modules.ModuleContext{
		DB:  db,
		JWT: jwt,
		Cfg: cfg,
		Redis: redis,
	}

	modules := []modules.ModuleInterface{
		auth.NewModule(ctx),
		user.NewModule(ctx),
	}

	api := r.Group("/api")

	for _, module := range modules {
		module.RegisterRoutes(api)
	}

	return &Application{
		router: r,
		config: cfg,
	}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr: a.config.GetPort(),
		Handler: a.router,
	}

	go func () {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
	}()

	log.Println("server running at port: ", a.config.GetPort())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	return srv.Shutdown(ctx)
}