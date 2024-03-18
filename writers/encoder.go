package writers

import (
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder() zapcore.Encoder {
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		_, file, line, ok := runtime.Caller(6)
		if ok {
			path := fmt.Sprintf("%s:%d", file, line)
			path = trimmedPath(path)
			enc.AppendString("[" + path + "]")
		} else {
			enc.AppendString("[" + caller.TrimmedPath() + "]")
		}

	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "caller_line"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
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
