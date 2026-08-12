package jobs

import (
	"context"
	"log/slog"
	"time"
)

func Run(ctx context.Context, log *slog.Logger, j Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("cronjob: panic recovered", "job", j.Name(), "panic", r)
		}
	}()

	start := time.Now()
	log.Info("cronjob: started", "job", j.Name())

	if err := j.Run(ctx, log); err != nil {
		log.Error("cronjob: failed", "job", j.Name(), "error", err, "duration", time.Since(start).String())
		return
	}

	log.Info("cronjob: finished", "job", j.Name(), "duration", time.Since(start).String())
}
