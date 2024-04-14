package sender

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/golang/protobuf/proto"
	"log"
	"sort"
	"strconv"
)

type ALiSLSConfig struct {
	DNS             string
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string // RAM用户角色的临时安全令牌，值为空表示不使用临时安全令牌。
	ProjectName     string
	LogStoreName    string
	Topic           string
	SampleInterval  int32
}

type ALiSLSSender struct {
	config   *ALiSLSConfig
	producer *producer.Producer
}

func NewALiSLSSender(config *ALiSLSConfig) *ALiSLSSender {
	if config == nil {
		log.Panic("ALiSLSSender's configuration is empty")
	}

	producerConfig := producer.GetDefaultProducerConfig()
	// 日志服务的服务入口。此处以杭州为例，其它地域请根据实际情况填写。
	producerConfig.Endpoint = config.DNS
	// 本示例从环境变量中获取AccessKey ID和AccessKey Secret。
	producerConfig.AccessKeyID = config.AccessKeyId
	producerConfig.AccessKeySecret = config.AccessKeySecret
	producerInstance := producer.InitProducer(producerConfig)

	return &ALiSLSSender{
		config:   config,
		producer: producerInstance,
	}
}

func (s *ALiSLSSender) Start() {
	s.producer.Start()
}

func (s *ALiSLSSender) Close() {
	s.producer.SafeClose()
}

func (s *ALiSLSSender) SendMetric(name string, tags map[string]string, value float64, ts int64, source string) error {
	project := s.config.ProjectName
	logstore := s.config.LogStoreName
	topic := s.config.Topic
	metric := buildLogItem(name, tags, value, ts)
	//log := buildLogItem("test_metric", map[string]string{"test_k": "test_v"}, value)
	//err := s.producer.SendLog(project, logstore, topic, "127.0.0.1", log)
	err := s.producer.SendLog(project, logstore, topic, source, metric)
	if err != nil {
		return err
	}
	return nil
}

func buildLogItem(metricName string, labels map[string]string, value float64, ts int64) *sls.Log {

	//now := uint32(time.Now().Unix())
	now := uint32(ts)
	log := &sls.Log{Time: proto.Uint32(now)}

	var contents []*sls.LogContent
	contents = append(contents, &sls.LogContent{
		Key:   proto.String("__time_nano__"),
		Value: proto.String(strconv.FormatInt(int64(now), 10)),
	})
	contents = append(contents, &sls.LogContent{
		Key:   proto.String("__name__"),
		Value: proto.String(metricName),
	})
	contents = append(contents, &sls.LogContent{
		Key:   proto.String("__value__"),
		Value: proto.String(strconv.FormatFloat(value, 'f', 6, 64)),
	})

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	labelsStr := ""
	for i, k := range keys {
		labelsStr += k
		labelsStr += "#$#"
		labelsStr += labels[k]
		if i < len(keys)-1 {
			labelsStr += "|"
		}
	}

	contents = append(contents, &sls.LogContent{Key: proto.String("__labels__"), Value: proto.String(labelsStr)})
	log.Contents = contents
	return log

}
