package log

import (
	"context"
)

type Log interface {
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, a ...interface{})

	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, a ...interface{})

	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, a ...interface{})

	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, format string, a ...interface{})

	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, format string, a ...interface{})

	SetLevel(level string)
}

var (
	logger Log
)

func GetInstance() Log {
	return logger
}

func InitLogger(l *Option) {
	logger = NewLogger(l)
}

func Debug(ctx context.Context, msg string) {
	logger.Debug(ctx, msg)
}

func Debugf(ctx context.Context, format string, a ...interface{}) {
	logger.Debugf(ctx, format, a...)
}

func Info(ctx context.Context, msg string) {
	logger.Info(ctx, msg)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	logger.Infof(ctx, format, a...)
}

func Warn(ctx context.Context, msg string) {
	logger.Warn(ctx, msg)
}

func Warnf(ctx context.Context, format string, a ...interface{}) {
	logger.Warnf(ctx, format, a...)
}

func Error(ctx context.Context, msg string) {
	logger.Error(ctx, msg)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	logger.Errorf(ctx, format, a...)
}

func Fatal(ctx context.Context, msg string) {
	logger.Fatal(ctx, msg)
}

func Fatalf(ctx context.Context, format string, a ...interface{}) {
	logger.Fatalf(ctx, format, a...)
}

const (
	// RequestID Request id key
	RequestID string = "request_id"
)

func SetLevel(level string) {
	logger.SetLevel(level)
}
