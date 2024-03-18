package saywo_logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var logger *SwLogger

func Register(opts ...LogOption) {
	logger = new(SwLogger)
	for _, opt := range opts {
		opt.apply(logger)
	}
	consoleInfoLog().apply(logger)
	consoleErrorLog().apply(logger)

	zlogger := zap.New(zapcore.NewTee(logger.cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.DPanicLevel))
	logger.l = zlogger
}

func Debug(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Debug(msg, zapFields...)
}

func Info(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Info(msg, zapFields...)
}

func Error(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Error(msg, zapFields...)
}

func Warn(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Warn(msg, zapFields...)
}

// 调用会触发panic和fatal
func Panic(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Panic(msg, zapFields...)
}

func DPanic(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.DPanic(msg, zapFields...)
}

func Fatal(msg string, fields map[string]any) {
	zapFields := getZapFileds(fields)
	logger.l.Fatal(msg, zapFields...)
}

func getZapFileds(fields map[string]any) []zap.Field {
	list := make([]zap.Field, 0)
	for key, value := range fields {
		list = append(list, zap.Any(key, value))
	}
	return list
}

func trimmedPath(path string) string {
	idx := len(path)
	for i := 0; i < 3; i++ {
		idx = strings.LastIndexByte(path[:idx], '/')
		if idx == -1 {
			return path
		}
	}
	file := path[idx+1:]
	return file
}

/* 适配旧项目*/

// InfoByArgs 通过参数输出日志
func InfoByArgs(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	logger.l.Info(s)
}

// ErrorByArgs 通过参数输出错误日志
func ErrorByArgs(format string, args ...interface{}) {
	logger.l.Error(fmt.Sprintf(format, args...))
}

// WarnByArgs 通过参数输出警告日志
func WarnByArgs(format string, args ...interface{}) {
	logger.l.Warn(fmt.Sprintf(format, args...))
}

// DebugByArgs 通过参数输出debug日志
func DebugByArgs(format string, args ...interface{}) {
	logger.l.Debug(fmt.Sprintf(format, args))
}

func ErrorKv(key string, value string, format string, args ...interface{}) {
	logger.l.Error(fmt.Sprintf(format, args), zap.String(key, value))

}

func InfoKv(key string, value string, format string, args ...interface{}) {
	logger.l.Info(fmt.Sprintf(format, args), zap.String(key, value))
}

func WarnKv(key string, value string, format string, args ...interface{}) {
	logger.l.Warn(fmt.Sprintf(format, args), zap.String(key, value))
}
