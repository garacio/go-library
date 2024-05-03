package log

import (
	"context"
	"log/slog"
	"os"
)

const LevelFatal slog.Level = 12

func Initialize(logger *slog.Logger) {
	slog.SetDefault(logger)
}

var globalLogLevel = new(slog.LevelVar)

// SetLevel ...
func SetLevel(level slog.Level) {
	globalLogLevel.Set(level)
}

func SetLevelFromString(level string) {
	switch level {
	case "debug":
		SetLevel(slog.LevelDebug)
	case "info":
		SetLevel(slog.LevelInfo)
	case "warn":
		SetLevel(slog.LevelWarn)
	case "error":
		SetLevel(slog.LevelError)
	case "fatal":
		SetLevel(LevelFatal)
	default:
		SetLevel(slog.LevelInfo)
	}
}

func Logger() *slog.Logger {
	output := os.Stdout
	logger := slog.New(
		&ContextHandler{
			Handler: PlainTextHandler{
				TextHandler: slog.NewTextHandler(output, &slog.HandlerOptions{
					Level: globalLogLevel,
				}),
				output: output,
			},
		},
	)

	return logger
}

// Info logs an info message
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error logs an error
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Fatal logs a fatal message
func Fatal(msg string, args ...any) {
	FatalCtx(context.Background(), msg, args...)
}

// InfoCtx logs an info message with context
func InfoCtx(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// WarnCtx logs a warn message with context
func WarnCtx(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// ErrorCtx logs an error message with context
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// DebugCtx logs a debug message with context
func DebugCtx(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// FatalCtx logs a fatal message with context
func FatalCtx(ctx context.Context, msg string, args ...any) {
	slog.Log(ctx, LevelFatal, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}
