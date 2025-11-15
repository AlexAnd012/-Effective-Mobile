package logging

import (
	"log/slog"
	"os"
	"strings"
)

func New(level string) *slog.Logger {
	var lvl slog.Level
	// Строке уровня возвращаем готовый *slog. Logger
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	// Создаём JSON handler, который пишет структурированные логи в stdout
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(h)
}
