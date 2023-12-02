package logsetup

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(opts ...LoggerOpt) *slog.Logger {
	cfg := &loggerOpts{
		level: slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.mode == LogModeDisabled {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.level,
	}
	if cfg.mode == LogModeJson {
		return slog.New(slog.NewJSONHandler(os.Stderr, slogOpts))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, slogOpts))
}
