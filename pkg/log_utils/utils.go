package log_utils

import "log/slog"

func MapSlogLevel(str string) slog.Level {
	switch str {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}

	return slog.LevelInfo
}
