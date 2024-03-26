package main

import (
	"os"

	"github.com/Edward-Alphonse/saywo_logs"
	"github.com/Edward-Alphonse/saywo_logs/writers"
	"github.com/pkg/errors"
)

type User struct {
	A int
	B *User
}

func main() {

	config := &writers.ALiSLSConfig{
		DNS:             "cn-wuhan-lr.log.aliyuncs.com",
		AccessKeyId:     os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"),
		Project:         "hzc-test-1",
		LogStore:        "aliyun-test-logstore",
		Topic:           "test",
		Level:           writers.DebugLevel,
	}
	fileConfig := &writers.FileConfig{
		Path:       writers.LogPath,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Level:      writers.DebugLevel,
	}
	saywo_logs.Register(saywo_logs.FileLog(fileConfig), saywo_logs.ALiSLS(config))
	saywo_logs.Debug("这是一个Debug")
	saywo_logs.Info("这是一个Info", []saywo_logs.Field{
		{"config": config},
		{"test": "test"},
	}...)
	saywo_logs.Warn("这是一个Warn")
	err := errors.New("1234")
	err = errors.Wrap(err, "FinishedCountStorage get value failed")

	user := User{
		A: 10,
		B: &User{
			A: 7,
		},
	}
	saywo_logs.Error("这是一个报错", saywo_logs.Field{
		"user_id":    1234456789,
		"article_id": 345645678900123456,
		"error":      err.Error(),
		"user":       &user,
	})
}
