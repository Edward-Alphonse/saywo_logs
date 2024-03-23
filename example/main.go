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
		ProjectName:     "hzc-test-1",
		LogStoreName:    "aliyun-test-logstore",
		Topic:           "test",
	}
	saywo_logs.Register(saywo_logs.FileLog(nil), saywo_logs.ALiSLS(config))
	saywo_logs.Info("这是一个Info", []saywo_logs.Field{
		{"config": config},
		{"test": "test"},
	}...)

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
	saywo_logs.Debug("这是一个debug")
}
