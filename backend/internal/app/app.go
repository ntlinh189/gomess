package app

import (
	"context"
	"gomess/internal/config"
	"gomess/internal/database"
	"gomess/internal/logger"
	"gomess/internal/middleware"
	"gomess/internal/modules"
	"gomess/internal/modules/auth"
	"gomess/internal/modules/friend"
	"gomess/internal/modules/message"
	"gomess/internal/modules/upload"
	"gomess/internal/modules/user"
	"gomess/internal/redis"
	"gomess/internal/ws"
	"gomess/pkg/jwt"
	"gomess/pkg/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	router *gin.Engine
	config config.ConfigInterface
	log    *slog.Logger
}

func NewApplication(cfg config.ConfigInterface) (*Application, error) {
	appLogger := logger.New(cfg.IsProduction())
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.WithLogger(appLogger))
	r.Use(middleware.Logging(appLogger))
	r.Use(middleware.Recovery(appLogger))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(cors.New(cors.Config{
		/// TODO: Change to client url
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: cfg.AllowCredentials(),
		MaxAge:           12 * time.Hour,
	}))

	db, err := database.NewMySql(cfg)
	if err != nil {
		return nil, err
	}

	jwt := jwt.NewJWT(cfg.GetJWTSecret())

	redis := redis.NewRedis(cfg.GetRedisAddr())

	ctx := &modules.ModuleContext{
		DB:    db,
		JWT:   jwt,
		Cfg:   cfg,
		Redis: redis,
	}

	storageClient, err := storage.NewStorage(
		cfg.GetMinioEndpoint(),
		cfg.GetMinioAccessKey(),
		cfg.GetMinioSecretKey(),
		cfg.GetMinioBucket(),
		cfg.MinioUseSSL(),
	)
	if err != nil {
		return nil, err
	}

	userRepo := user.NewRepository(db)
	friendRepo := friend.NewRepository(db)
	messageRepo := message.NewRepository(db)

	hub := ws.NewHub()

	modules := []modules.ModuleInterface{
		auth.NewModule(ctx),
		user.NewModule(ctx, userRepo),
		friend.NewModule(ctx, friendRepo, userRepo, messageRepo),
		message.NewModule(ctx, friendRepo, storageClient, hub),
		upload.NewModule(ctx, storageClient),
		ws.NewModule(ctx, hub),
	}

	api := r.Group("/api")

	for _, module := range modules {
		module.RegisterRoutes(api)
	}

	return &Application{
		router: r,
		config: cfg,
		log:    appLogger,
	}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.GetPort(),
		Handler: a.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("server failed to start", "error", err)
			return
		}
	}()

	a.log.Info("server running", "port", a.config.GetPort())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	a.log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	return srv.Shutdown(ctx)
}
