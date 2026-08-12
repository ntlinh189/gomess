package jobs

import (
	"context"
	"log/slog"
)

type Job interface {
	Name() string
	Schedule() string
	Run(ctx context.Context, log *slog.Logger) error
}