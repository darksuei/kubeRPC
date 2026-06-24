package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Init() {
	var level slog.Level
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default: // DEBUG or unset
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}
