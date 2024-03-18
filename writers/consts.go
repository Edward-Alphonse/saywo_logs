package writers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel  Level = zap.DebugLevel
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel // debug环境下使用
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel // 会调用os.Exit(1) 退出
)

const logPath = "./run"
