package logger

import (
	"log/slog"
	"os"
)

func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Error(msg string, err error) {
	slog.Error(msg, "error", err)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}
