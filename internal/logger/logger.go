package logger

import "log/slog"

var (
	logger *slog.Logger
)

func GetLogger() *slog.Logger {
	if logger == nil {
		logger = slog.With("app", "mineserver")
	}
	return logger
}
