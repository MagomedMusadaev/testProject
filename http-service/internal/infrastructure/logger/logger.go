package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger представляет собой обертку для slog
type Logger struct {
	logger *slog.Logger
}

// NewLogger создает новый экземпляр логгера с уровнем логирования Info
func NewLogger() *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		logger: slog.New(handler),
	}
}

// Info логирует информационное сообщение
// msg - информационное сообщение
// args - дополнительные параметры для логирования
func (l *Logger) Info(msg string, args ...any) {
	l.logger.InfoContext(context.Background(), msg, append([]any{}, args...)...)
}

// Error логирует сообщение об ошибке
// op - название операции
// msg - сообщение об ошибке
// err - объект ошибки
// args - дополнительные параметры для логирования
func (l *Logger) Error(op string, msg string, err error, args ...any) {
	l.logger.ErrorContext(context.Background(), msg, append([]any{"op", op, "error", err}, args...)...)
}
