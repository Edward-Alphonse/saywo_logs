package writers

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
1. 保留的日志是否被删除取决于MaxBackups和MaxAge，满足二者中的一个条件就会触发日志文件的删除
2. 删除日志逻辑先判断MaxBackups，旧日志数超过MaxBackups，则删除一部分日志
3. 删除日志逻辑然后判断MaxAge，超过存储天数的日志将被删除
4. 经过两轮删除后，留存的日志一定满足MaxBackups和MaxAge的条件
4. MaxBackups 和 MaxAge 都是 0，则不会删除任何旧日志文件
*/

type FileConfig struct {
	Path       string //日志文件路径
	MaxSize    int    // 日志单个文件的大小，以兆字节为单位
	MaxBackups int    // 已被分割存储的日志文件最多的留存个数，单位是个
	MaxAge     int    // 已被分割存储的日志文件最大的留存时间，单位是天
	Compress   bool   //指定被分割之后的文件是否要压缩
	Level      Level  //日志等级输入控制
}

func NewFileCore(config *FileConfig) zapcore.Core {
	if config == nil {
		config = getDefaultConfig()
	}
	writeSyncer := getWriterSyncer(config)
	encoder := getEncoder()
	lvl := GetZapCoreLevel(config.Level)
	core := zapcore.NewCore(encoder, writeSyncer, lvl)
	return core
}

func getDefaultConfig() *FileConfig {
	config := &FileConfig{
		Path:       LogPath,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Level:      InfoLevel,
	}
	return config
}

func getWriterSyncer(config *FileConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Path + getLogFileName(),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogFileName() string {
	return fmt.Sprintf("_log.log")
}
