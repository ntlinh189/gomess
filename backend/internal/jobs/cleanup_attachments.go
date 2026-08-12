package jobs

import (
	"context"
	"fmt"
	"gomess/internal/modules/message"
	"gomess/pkg/storage"
	"log/slog"
	"time"
)

const attachmentGracePeriod = 24 * time.Hour

type CleanupAttachmentsJob struct {
	repo     message.RepositoryInterface
	storage  storage.StorageInterface
	schedule string
}

func NewCleanupAttachmentsJob(
	repo message.RepositoryInterface,
	s storage.StorageInterface,
	intervalHours int,
) *CleanupAttachmentsJob {
	return &CleanupAttachmentsJob{
		repo:     repo,
		storage:  s,
		schedule: fmt.Sprintf("@every %dh", intervalHours),
	}
}

func (j *CleanupAttachmentsJob) Name() string     { return "cleanup-attachments" }
func (j *CleanupAttachmentsJob) Schedule() string { return j.schedule }

func (j *CleanupAttachmentsJob) Run(ctx context.Context, log *slog.Logger) error {
	referenced, err := j.repo.GetAllObjectKeys()
	if err != nil {
		return fmt.Errorf("list referenced object keys: %w", err)
	}

	referencedSet := make(map[string]struct{}, len(referenced))
	for _, key := range referenced {
		referencedSet[key] = struct{}{}
	}

	stored, err := j.storage.ListObjects(ctx, "attachments/")
	if err != nil {
		return fmt.Errorf("list storage objects: %w", err)
	}

	deleted := 0
	failed := 0

	for _, obj := range stored {
		if _, ok := referencedSet[obj.Key]; ok {
			continue
		}
		if time.Since(obj.LastModified) < attachmentGracePeriod {
			continue
		}
		if err := j.storage.DeleteObject(ctx, obj.Key); err != nil {
			failed++
			log.Error("cronjob: failed to delete orphaned object", "key", obj.Key, "error", err)
			continue
		}
		deleted++
	}

	log.Info("cronjob: cleanup-attachments summary",
		"scanned", len(stored),
		"deleted", deleted,
		"failed", failed,
	)

	return nil
}