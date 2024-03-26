package writers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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
	SecurityToken   string // RAM用户角色的临时安全令牌，值为空表示不使用临时安全令牌。
	Project         string
	LogStore        string
	Topic           string
	Level           Level
}

func NewALiLogCore(config *ALiSLSConfig) zapcore.Core {
	writeSyncer := NewALiSLSWriter(config)
	encoder := getEncoder()
	lvl := GetZapCoreLevel(config.Level)
	return zapcore.NewCore(encoder, writeSyncer, lvl)
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
	securityToken := config.SecurityToken
	endPoint := config.DNS

	// 创建日志服务Client。
	provider := sls.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, securityToken)
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
		val, err := getStringValue(value)
		if err != nil {
			log.Printf("ALiSLSWriter get string value failed, error: %v", err)
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
		Source: proto.String("203.0.113.10"), // todo：使用服务器IP
		Logs:   []*sls.Log{logMsg},
	}

	err = w.client.PutLogs(w.config.Project, w.config.LogStore, loggroup)
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

func getStringValue(value any) (string, error) {
	vt := reflect.TypeOf(value)
	numberType := reflect.TypeOf(json.Number(""))
	var val string
	switch vt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = fmt.Sprintf("%d", value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = fmt.Sprintf("%d", value)
	case reflect.String:
		if vt.AssignableTo(numberType) {
			value = value.(json.Number).String()
		}
		val = value.(string)
	case reflect.Float32, reflect.Float64:
		val = fmt.Sprintf("%f", value)
	case reflect.Bool:
		val = strconv.FormatBool(value.(bool))
	case reflect.Array, reflect.Map, reflect.Struct:
		bytes, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("ALiSLSWriter write failed, err: %v", err)
		}
		val = string(bytes)
	default:
		return "", fmt.Errorf("Unknown type: %v", vt)
	}
	return val, nil
}
