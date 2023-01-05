package logger

import "golang.org/x/exp/slog"

var Default Logger = slog.Default()

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, err error, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}
