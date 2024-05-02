package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/kosalnik/gmarket/internal/config"
)

var loggerInstance *slog.Logger

var levelMap = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func InitLogger(cfg config.Logger) error {
	level, ok := levelMap[strings.ToUpper(cfg.Level)]
	if !ok {
		return errors.New("unknown name")
	}
	loggerInstance = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	return nil
}

func Debug(msg string, args ...any) {
	loggerInstance.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	loggerInstance.Warn(msg, args...)
}

func Info(msg string, args ...any) {
	loggerInstance.Info(msg, args...)
}

func Error(msg string, args ...any) {
	loggerInstance.Error(msg, args...)
}
