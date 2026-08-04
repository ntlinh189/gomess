package logger

import (
	"log/slog"
	"os"
)

func New(isProduction bool) *slog.Logger {
	level := slog.LevelDebug
	if isProduction {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if isProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}