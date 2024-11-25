package logger

import (
	"context"
	"dockside/app/util/must"
	"log/slog"
)

func Error(message string, err error) {
	slog.Error(message, slog.String("error", err.Error()))
}

func ErrorWithFields(message string, err error, fields map[string]any) {
	if err != nil {
		fields["error"] = err.Error()
	}
	MessageWithFields(message, slog.LevelError, fields)
}

func Info(message string) {
	slog.Info(message)
}

func InfoWithFields(message string, fields map[string]any) {
	MessageWithFields(message, slog.LevelInfo, fields)
}

func Debug(message string) {
	slog.Debug(message)
}

func DebugWithFields(message string, fields map[string]any) {
	MessageWithFields(message, slog.LevelDebug, fields)
}

func Warn(message string) {
	slog.Warn(message)
}

func WarnWithFields(message string, fields map[string]any) {
	MessageWithFields(message, slog.LevelWarn, fields)
}

func MessageWithFields(message string, level slog.Level, fields map[string]any) {
	attrs := make([]slog.Attr, 0, len(fields))
	for s, a := range fields {
		attrs = append(attrs, slog.String(s, string(must.Serialize(a))))
	}
	slog.LogAttrs(context.Background(), level, message, attrs...)
}
