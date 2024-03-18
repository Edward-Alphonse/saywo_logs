package writers

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(lvl Level) bool

func NewConsoleCore(level Level, Lef LevelEnablerFunc) zapcore.Core {
	encoder := getEncoder()
	writerSyncer := os.Stdout
	if level < ErrorLevel {
		writerSyncer = os.Stderr
	}
	lv := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return Lef(level)
	})
	core := zapcore.NewCore(encoder, writerSyncer, lv)
	return core
}

// ErrorLevel级别以下日志输出到stdout
func InfoConsoleCore() zapcore.Core {
	infoCore := NewConsoleCore(InfoLevel, func(lvl Level) bool {
		return lvl < ErrorLevel
	})
	return infoCore
}

// ErrorLevel级别及以上日志输出到stderr
func ErrorConsoleCore() zapcore.Core {
	errCore := NewConsoleCore(ErrorLevel, func(lvl Level) bool {
		return lvl >= ErrorLevel
	})
	return errCore
}

func DefaultConsoleCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0)
	cores = append(cores, InfoConsoleCore())
	cores = append(cores, ErrorConsoleCore())
	return cores
}
