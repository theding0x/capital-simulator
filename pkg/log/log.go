// Package log provides a small wrapper around log/slog so every service in
// capital-simulator emits structured logs in a consistent shape.
package log

import (
	"log/slog"
	"os"
)

// New returns a configured *slog.Logger suitable for service entrypoints.
// The "service" field is attached to every record.
func New(service string) *slog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler).With(slog.String("service", service))
}

// SetDefault installs the given logger as the slog default.
func SetDefault(l *slog.Logger) {
	slog.SetDefault(l)
}

func parseLevel(s string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo
	}
	return l
}
