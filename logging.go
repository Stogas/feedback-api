package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Stogas/feedback-api/internal/config"
)

func initLogger(conf config.LogsConfig) {
	var handler slog.Handler

	level := slog.LevelInfo
	if conf.Debug {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	if conf.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
}

func getLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
