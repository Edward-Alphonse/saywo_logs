package writers

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(lvl zapcore.Level) bool

func NewConsoleCore(level Level, Lef LevelEnablerFunc) zapcore.Core {
	encoder := getEncoder()
	writerSyncer := os.Stdout
	lvl := GetZapCoreLevel(level)
	if lvl >= zapcore.ErrorLevel {
		writerSyncer = os.Stderr
	}
	lv := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return Lef(level)
	})
	core := zapcore.NewCore(encoder, writerSyncer, lv)
	return core
}

// Debug级别日志输出控制
func DebugConsoleCore() zapcore.Core {
	infoCore := NewConsoleCore(DebugLevel, func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	return infoCore
}

// ErrorLevel级别以下日志输出到stdout
func InfoConsoleCore() zapcore.Core {
	infoCore := NewConsoleCore(InfoLevel, func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl > zapcore.DebugLevel
	})
	return infoCore
}

// ErrorLevel级别及以上日志输出到stderr
func ErrorConsoleCore() zapcore.Core {
	errCore := NewConsoleCore(ErrorLevel, func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	return errCore
}
