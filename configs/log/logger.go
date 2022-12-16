package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

func (l *Logger) Debug(ctx context.Context, msg string) {
	evt := l.Logger.Debug()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(msg)
}

func (l *Logger) Debugf(ctx context.Context, format string, a ...interface{}) {
	evt := l.Logger.Debug()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(fmt.Sprintf(format, a...))
}

func (l *Logger) Info(ctx context.Context, msg string) {
	evt := l.Logger.Info()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msg(msg)
}

func (l *Logger) Infof(ctx context.Context, format string, a ...interface{}) {
	evt := l.Logger.Info()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(fmt.Sprintf(format, a...))
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	evt := l.Logger.Warn()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msg(msg)
}

func (l *Logger) Warnf(ctx context.Context, format string, a ...interface{}) {
	evt := l.Logger.Warn()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(fmt.Sprintf(format, a...))
}

func (l *Logger) Error(ctx context.Context, msg string) {
	evt := l.Logger.Error()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msg(msg)
}

func (l *Logger) Errorf(ctx context.Context, format string, a ...interface{}) {
	evt := l.Logger.Error()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(fmt.Sprintf(format, a...))
}

func (l *Logger) Fatal(ctx context.Context, msg string) {
	evt := l.Logger.Fatal()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msg(msg)
}

func (l *Logger) Fatalf(ctx context.Context, format string, a ...interface{}) {
	evt := l.Logger.Fatal()
	for k, v := range withCtx(ctx) {
		evt = evt.Str(" | "+k, v)
	}
	evt.Msgf(fmt.Sprintf(format, a...))
}

func (l *Logger) SetLevel(level string) {
	Level, err := zerolog.ParseLevel(level)
	if err != nil {
		l.Logger.Error().Msgf("setLevel err:%v", err)
		return
	}
	l.Logger.Level(Level)
}

func withCtx(ctx context.Context) (kv map[string]string) {
	meta, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		meta, ok = metadata.FromIncomingContext(ctx)
	}
	kv = map[string]string{}
	if !ok {
		requestID, _ := ctx.Value(RequestID).(string)
		kv[RequestID] = requestID
		return
	}

	getVal := func(val []string) string {
		if len(val) == 0 {
			return ""
		}
		return val[0]
	}

	kv[RequestID] = getVal(meta.Get(RequestID))
	return
}

func NewLogger(l *Option) Log {
	if l.MaxFileSize == "" {
		l.MaxFileSize = "500M"
	}

	if l.Level == "" {
		l.Level = "info"
	}

	if l.RotateDuration == "" {
		l.RotateDuration = "1h"
	}

	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat

	// 创建log目录
	logDir := l.DirPath
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("init log Mkdir failed, err:", err)
		return nil
	}
	fileName := logDir + time.Now().Format("2006-01-02") + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s;", i)
	}
	logger := zerolog.New(zerolog.MultiLevelWriter(consoleWriter, logFile)).With().Timestamp().Logger()
	return &Logger{
		Logger: &logger,
	}
}
