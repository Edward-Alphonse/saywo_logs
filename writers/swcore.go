package writers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 直接复制的zapCore的ioCore的实现，以后自定义再说，目前可以平替zapCore.NewCore
// github上另外的实现，https://github.com/yeyudekuangxiang/zap-aliyun-log
type SwCore struct {
	level Level
	enc   zapcore.Encoder
	out   zapcore.WriteSyncer
}

func NewSwCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, level Level) zapcore.Core {
	return &SwCore{
		enc:   enc,
		out:   ws,
		level: level,
	}
}

func (c *SwCore) Enabled(level Level) bool {
	return level >= c.level
}

// With adds structured context to the Core.
func (c *SwCore) With(fields []zap.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

func (c *SwCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *SwCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > ErrorLevel {
		_ = c.Sync()
	}
	return nil
}

func (c *SwCore) Sync() error {
	return c.out.Sync()
}

func (c *SwCore) clone() *SwCore {
	newCore := &SwCore{
		level: c.level,
		enc:   c.enc.Clone(),
	}
	return newCore
}

func addFields(enc zapcore.Encoder, fields []zap.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
