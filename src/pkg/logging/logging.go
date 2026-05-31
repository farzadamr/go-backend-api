package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/farzadamr/go-backend-api/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

const (
	TraceIdKey contextKey = "trace_id"
	UserIdKey  contextKey = "user_id"
)

func Init(cfg *config.Config) error {
	var level slog.Level
	switch cfg.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			if a.Key == slog.LevelKey {
				a.Key = "severity"
			}
			return a
		},
	}
	// Determine output writers
	var writers []io.Writer
	writers = append(writers, os.Stdout) // always log to stdout

	if cfg.Logger.LogFile != "" {
		// Ensure directory exists
		dir := filepath.Dir(cfg.Logger.LogFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		// Lumberjack logger
		lumber := &lumberjack.Logger{
			Filename:   cfg.Logger.LogFile,
			MaxSize:    cfg.Logger.MaxSize,
			MaxBackups: cfg.Logger.MaxBackups,
			MaxAge:     cfg.Logger.MaxAge,
			Compress:   cfg.Logger.Compress,
		}
		writers = append(writers, lumber)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create handler (JSON or Text)
	var handler slog.Handler
	if cfg.Logger.JSONFormat {
		handler = slog.NewJSONHandler(multiWriter, opts)
	} else {
		handler = slog.NewTextHandler(multiWriter, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

// WithTraceID returns a new context with trace_id value
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIdKey, traceID)
}

// WithUserID returns a new context with user_id value
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIdKey, userID)
}

// getTraceIDFromContext extracts trace_id from context
func getTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(TraceIdKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// getUserIDFromContext extracts user_id from context
func getUserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(UserIdKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// buildContextAttrs builds slog attributes from context values
func buildContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if userID := getUserIDFromContext(ctx); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	return attrs
}

// --- Convenience functions that accept context ---

// DebugContext logs debug message with context fields
func DebugContext(ctx context.Context, msg string, args ...any) {
	attrs := buildContextAttrs(ctx)
	if len(attrs) > 0 {
		args = append(attrsToArgs(attrs), args...)
	}
	slog.Debug(msg, args...)
}

// InfoContext logs info message with context fields
func InfoContext(ctx context.Context, msg string, args ...any) {
	attrs := buildContextAttrs(ctx)
	if len(attrs) > 0 {
		args = append(attrsToArgs(attrs), args...)
	}
	slog.Info(msg, args...)
}

// WarnContext logs warn message with context fields
func WarnContext(ctx context.Context, msg string, args ...any) {
	attrs := buildContextAttrs(ctx)
	if len(attrs) > 0 {
		args = append(attrsToArgs(attrs), args...)
	}
	slog.Warn(msg, args...)
}

// ErrorContext logs error message with context fields
func ErrorContext(ctx context.Context, msg string, args ...any) {
	attrs := buildContextAttrs(ctx)
	if len(attrs) > 0 {
		args = append(attrsToArgs(attrs), args...)
	}
	slog.Error(msg, args...)
}

// attrsToArgs converts []slog.Attr to slog.Attr variadic args
func attrsToArgs(attrs []slog.Attr) []any {
	res := make([]any, len(attrs))
	for i, a := range attrs {
		res[i] = a
	}
	return res
}

// --- Legacy functions (without context) for places where context is not available
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
