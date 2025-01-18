package utils

import (
	"log/slog"
	"os"
)

func JsonLog(msg string, level string, data interface{}) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	switch level {
	case "info":
		logger.Info(msg, slog.Any("data", data))
	case "error":
		logger.Error(msg, slog.Any("data", data))
	case "warn":
		logger.Warn(msg, slog.Any("data", data))
	case "debug":
		logger.Debug(msg, slog.Any("data", data))
	case "fatal":
		logger.Error(msg, slog.Any("data", data))
		os.Exit(1)
	default:
		logger.Info(msg, slog.Any("data", data))
	}
}
