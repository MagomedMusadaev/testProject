package logger

import (
	"log/slog"
	"os"
)

// Logger представляет собой обертку вокруг slog
type Logger struct {
	*slog.Logger
}

// NewLogger создает новый экземпляр логгера с JSON форматированием
func NewLogger() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Error логирует сообщение об ошибке с контекстом операции
// op - название операции
// msg - сообщение об ошибке
// err - объект ошибки
func (l *Logger) Error(op string, msg string, err error) {
	l.Logger.Error(msg,
		slog.String("op", op),
		slog.String("error", err.Error()),
	)
}

// Info логирует информационное сообщение с контекстом операции
// msg - информационное сообщение
// args - дополнительные параметры для логирования
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg,
		slog.Any("details", args),
	)
}
