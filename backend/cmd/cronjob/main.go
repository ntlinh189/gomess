package main

import (
	"context"
	"gomess/internal/config"
	"gomess/internal/database"
	"gomess/internal/jobs"
	"gomess/internal/logger"
	"gomess/internal/modules/message"
	"gomess/pkg/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

const cleanupIntervalHours = 6 // Hours

func main() {
	cfg := config.NewConfig()
	log := logger.New(cfg.IsProduction())

	db, err := database.NewMySql(cfg)
	if err != nil {
		log.Error("cronjob: db connect failed", "error", err)
		os.Exit(1)
	}

	storageClient, err := storage.NewStorage(
		cfg.GetMinioEndpoint(), cfg.GetMinioPublicEndpoint(), cfg.GetMinioAccessKey(),
		cfg.GetMinioSecretKey(), cfg.GetMinioBucket(),
	)
	if err != nil {
		log.Error("cronjob: storage connect failed", "error", err)
		os.Exit(1)
	}

	messageRepo := message.NewRepository(db)

	registeredJobs := []jobs.Job{
		jobs.NewCleanupAttachmentsJob(messageRepo, storageClient, cleanupIntervalHours),
	}

	c := cron.New()

	for _, j := range registeredJobs {
		job := j
		if _, err := c.AddFunc(job.Schedule(), func() {
			jobs.Run(context.Background(), log, job)
		}); err != nil {
			log.Error("cronjob: register failed", "job", job.Name(), "error", err)
			os.Exit(1)
		}
		log.Info("cronjob: registered", "job", job.Name(), "schedule", job.Schedule())
	}

	c.Start()
	log.Info("cronjob: scheduler running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("cronjob: shutting down")
	stopCtx := c.Stop()
	<-stopCtx.Done()
}
