// Package log provides a small wrapper around log/slog so every service in
// capital-simulator emits structured logs in a consistent shape.
package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
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

// FromContext returns a logger associated with ctx if one was attached via
// WithContext, otherwise the slog default.
func FromContext(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok && v != nil {
		return v
	}
	return slog.Default()
}

// WithContext returns a copy of ctx with l attached.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

type loggerKey struct{}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "", "info":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}
