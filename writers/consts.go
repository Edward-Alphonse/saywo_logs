package writers

import (
	"go.uber.org/zap/zapcore"
)

type Level = string

const (
	DebugLevel  Level = "debug"
	InfoLevel   Level = "info"
	WarnLevel   Level = "warn"
	ErrorLevel  Level = "error"
	DPanicLevel Level = "dpanic" // debug环境下使用
	PanicLevel  Level = "panic"
	FatalLevel  Level = "fatal"
)

var dict = map[Level]zapcore.Level{
	DebugLevel:  zapcore.DebugLevel,
	InfoLevel:   zapcore.InfoLevel,
	WarnLevel:   zapcore.WarnLevel,
	ErrorLevel:  zapcore.ErrorLevel,
	DPanicLevel: zapcore.DPanicLevel, // debug环境下使用
	PanicLevel:  zapcore.PanicLevel,
	FatalLevel:  zapcore.FatalLevel, // 会调用os.Exit(1) 退出
}

func GetZapCoreLevel(lvl Level) zapcore.Level {
	res, ok := dict[lvl]
	if !ok {
		return zapcore.InfoLevel
	}
	return res
}

const LogPath = "./run"
