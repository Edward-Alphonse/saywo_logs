package saywo_logs

import (
	"github.com/Edward-Alphonse/saywo_logs/writers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogOption interface {
	apply(logger *SwLogger)
}

type LogOptionFunc func() zapcore.Core

func (f LogOptionFunc) apply(logger *SwLogger) {
	cores := logger.cores
	if len(cores) == 0 {
		cores = make([]zapcore.Core, 0)
	}
	logger.cores = append(cores, f())
}

type SwLogger struct {
	cores []zapcore.Core
	l     *zap.Logger
}

func consoleInfoLog() LogOption {
	function := func() zapcore.Core {
		return writers.InfoConsoleCore()
	}
	return LogOptionFunc(function)
}

func consoleErrorLog() LogOption {
	function := func() zapcore.Core {
		return writers.ErrorConsoleCore()
	}
	return LogOptionFunc(function)
}

func FileLog(config *writers.FileConfig) LogOption {
	return LogOptionFunc(func() zapcore.Core {
		return writers.NewFileCore(config)
	})
}

func ALiSLS(config *writers.ALiSLSConfig) LogOption {
	return LogOptionFunc(func() zapcore.Core {
		return writers.NewALiLogCore(config)
	})
}
