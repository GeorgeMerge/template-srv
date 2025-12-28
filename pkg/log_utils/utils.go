package log_utils

import "log/slog"

func MapSlogLevel(str string) slog.Level {
	switch str {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelInfo
}
