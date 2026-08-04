package main

import (
	"context"
	"gomess/internal/config"
	"gomess/internal/database"
	"gomess/internal/logger"
	"gomess/internal/modules/message"
	"gomess/pkg/storage"
	"log/slog"
	"os"
	"time"
)

const gracePeriod = 24 * time.Hour

func main() {
	cfg := config.NewConfig()
	log := logger.New(cfg.IsProduction())

	db, err := database.NewMySql(cfg)
	if err != nil {
		log.Error("worker: failed to connect database", "error", err)
		os.Exit(1)
	}

	storageClient, err := storage.NewStorage(
		cfg.GetMinioEndpoint(),
		cfg.GetMinioAccessKey(),
		cfg.GetMinioSecretKey(),
		cfg.GetMinioBucket(),
		cfg.MinioUseSSL(),
	)
	if err != nil {
		log.Error("worker: failed to connect storage", "error", err)
		os.Exit(1)
	}

	messageRepo := message.NewRepository(db)

	interval := time.Duration(cfg.WorkerCleanupIntervalHours()) * time.Hour
	log.Info("worker started", "cleanup_interval", interval.String())

	for {
		cleanupOrphanedAttachments(context.Background(), log, messageRepo, storageClient)
		time.Sleep(interval)
	}
}

func cleanupOrphanedAttachments(
	ctx context.Context,
	log *slog.Logger,
	repo message.RepositoryInterface,
	s storage.StorageInterface,
) {
	log.Info("cleanup job started")

	referenced, err := repo.GetAllObjectKeys()
	if err != nil {
		log.Error("cleanup: failed to list referenced object keys", "error", err)
		return
	}

	referencedSet := make(map[string]struct{}, len(referenced))
	for _, key := range referenced {
		referencedSet[key] = struct{}{}
	}

	stored, err := s.ListObjects(ctx, "attachments/")
	if err != nil {
		log.Error("cleanup: failed to list storage objects", "error", err)
		return
	}

	deleted := 0
	for _, obj := range stored {
		if _, ok := referencedSet[obj.Key]; ok {
			continue
		}

		if time.Since(obj.LastModified) < gracePeriod {
			continue
		}

		if err := s.DeleteObject(ctx, obj.Key); err != nil {
			log.Error("cleanup: failed to delete orphaned object", "object_key", obj.Key, "error", err)
			continue
		}

		deleted++
	}

	log.Info("cleanup job finished", "checked", len(stored), "deleted", deleted)
}