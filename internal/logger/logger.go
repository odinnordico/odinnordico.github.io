package logger

import (
	"log/slog"
)

var logger *slog.Logger

func init() {
	logger = slog.Default()
}

func SetupLogger(l *slog.Logger) {
	logger = l
}

func Logger() *slog.Logger {
	return logger
}
