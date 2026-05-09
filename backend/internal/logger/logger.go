package logger

import (
	"log/slog"
	"os"
)

func InitLogger(debug bool) {
	logger := slog.New(
		func() slog.Handler {
			if debug {
				return slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				)
			}
			return slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			)
		}(),
	)
	slog.SetDefault(logger)
}
