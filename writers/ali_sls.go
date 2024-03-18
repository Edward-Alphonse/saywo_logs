package writers

import (
	"log"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	simpleJson "github.com/bitly/go-simplejson"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap/zapcore"
)

type ALiSLSConfig struct {
	DNS             string
	AccessKeyId     string
	AccessKeySecret string
	ProjectName     string
	LogStoreName    string
	Topic           string
}

func NewALiLogCore(config *ALiSLSConfig) zapcore.Core {
	writeSyncer := NewALiSLSWriter(config)
	encoder := getEncoder()
	return zapcore.NewCore(encoder, writeSyncer, InfoLevel)
}

type ALiSLSWriter struct {
	config *ALiSLSConfig
	client sls.ClientInterface
}

func NewALiSLSWriter(config *ALiSLSConfig) *ALiSLSWriter {
	if config == nil {
		log.Panic("ALiSLSWriter's configuration is empty")
	}
	accessKeyId := config.AccessKeyId
	accessKeySecret := config.AccessKeySecret
	endPoint := config.DNS
	// RAM用户角色的临时安全令牌。此处取值为空，表示不使用临时安全令牌。
	SecurityToken := ""
	// 创建日志服务Client。
	provider := sls.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, SecurityToken)
	client := sls.CreateNormalInterfaceV2(endPoint, provider)

	writer := &ALiSLSWriter{
		config: config,
		client: client,
	}
	return writer
}

func (w *ALiSLSWriter) Write(p []byte) (n int, err error) {
	jsonObject, err := simpleJson.NewJson(p)
	if err != nil {
		log.Printf("ALiSLSWriter new json failed %v", err)
		return 0, err
	}

	m, err := jsonObject.Map()
	if err != nil {
		log.Printf("ALiSLSWriter new json failed %v", err)
		return 0, err
	}
	content := []*sls.LogContent{}
	for key, value := range m {
		val, ok := value.(string)
		if !ok {
			continue
		}
		content = append(content, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(val),
		})
	}
	logMsg := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}

	loggroup := &sls.LogGroup{
		Topic:  proto.String(w.config.Topic),
		Source: proto.String("203.0.113.10"),
		Logs:   []*sls.Log{logMsg},
	}

	err = w.client.PutLogs(w.config.ProjectName, w.config.LogStoreName, loggroup)
	if err != nil {
		log.Fatalf("PutLogs failed %v", err)
		os.Exit(1)
		return 0, err
	}
	log.Println("PutLogs success")
	return len(p), nil
}

func (w *ALiSLSWriter) Sync() error {
	return nil
}
