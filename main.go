package main

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"os"
	"time"
)

//	"fmt"
//	sls "github.com/aliyun/aliyun-log-go-sdk"
//	"github.com/aliyun/aliyun-log-go-sdk/producer"
//	"github.com/golang/protobuf/proto"
//	"github.com/rcrowley/go-metrics"
//	"os"
//	"os/signal"
//	"sort"
//	"strconv"
//	"sync"
//	"time"

//func buildLogItem(metricName string, labels map[string]string, value float64) *sls.Log {
//
//	now := uint32(time.Now().Unix())
//	log := &sls.Log{Time: proto.Uint32(now)}
//
//	var contents []*sls.LogContent
//	contents = append(contents, &sls.LogContent{
//		Key:   proto.String("__time_nano__"),
//		Value: proto.String(strconv.FormatInt(int64(now), 10)),
//	})
//	contents = append(contents, &sls.LogContent{
//		Key:   proto.String("__name__"),
//		Value: proto.String(metricName),
//	})
//	contents = append(contents, &sls.LogContent{
//		Key:   proto.String("__value__"),
//		Value: proto.String(strconv.FormatFloat(value, 'f', 6, 64)),
//	})
//
//	keys := make([]string, 0, len(labels))
//	for k := range labels {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//
//	labelsStr := ""
//	for i, k := range keys {
//		labelsStr += k
//		labelsStr += "#$#"
//		labelsStr += labels[k]
//		if i < len(keys)-1 {
//			labelsStr += "|"
//		}
//	}
//
//	contents = append(contents, &sls.LogContent{Key: proto.String("__labels__"), Value: proto.String(labelsStr)})
//	log.Contents = contents
//	return log
//
//}

//func main() {
//
//	config := &ALiSLSConfig{
//		DNS:             "cn-wuhan-lr.log.aliyuncs.com",
//		AccessKeyId:     os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"),
//		AccessKeySecret: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"),
//		ProjectName:     "hzc-test",
//		LogStoreName:    "metrics_test",
//		Topic:           "test",
//	}
//
//	project := config.ProjectName
//	logstore := config.LogStoreName
//	topic := config.Topic
//
//	producerConfig := producer.GetDefaultProducerConfig()
//	// 日志服务的服务入口。此处以杭州为例，其它地域请根据实际情况填写。
//	producerConfig.Endpoint = config.DNS
//	// 本示例从环境变量中获取AccessKey ID和AccessKey Secret。
//	producerConfig.AccessKeyID = config.AccessKeyId
//	producerConfig.AccessKeySecret = config.AccessKeySecret
//	producerInstance := producer.InitProducer(producerConfig)
//	ch := make(chan os.Signal)
//	signal.Notify(ch)
//	producerInstance.Start()
//	var m sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		m.Add(1)
//		go func() {
//			defer m.Done()
//			for i := 0; i < 1000; i++ {
//				// GenerateLog  is producer's function for generating SLS format logs
//				// GenerateLog has low performance, and native Log interface is the best choice for high performance.
//				log := buildLogItem("test_metric", map[string]string{"test_k": "test_v"}, float64(i))
//				err := producerInstance.SendLog(project, logstore, topic, "127.0.0.1", log)
//				if err != nil {
//					fmt.Println(err)
//				}
//			}
//		}()
//	}
//	m.Wait()
//	fmt.Println("Send completion")
//	if _, ok := <-ch; ok {
//		fmt.Println("Get the shutdown signal and start to shut down")
//		producerInstance.Close(60000)
//	}
//}

func main() {
	counter := metrics.NewCounter()
	metrics.Register("counter_name", counter)
	counter.Inc(47)

	g := metrics.NewGauge()
	err := metrics.Register("gauge_name", g)
	if err != nil {
		fmt.Printf("Register Gauge %v\n", err)
		return
	}
	g.Update(47)
	g.Value()

	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register("histogram_name", h)
	h.Update(47)
	for i := 0; i < 100; i++ {
		h.Update(int64(i))
	}

	m := metrics.NewMeter()
	metrics.Register("meter_name", m)
	m.Mark(47)
	for i := 0; i < 10; i++ {
		m.Mark(int64(1))
		time.Sleep(time.Millisecond * 1)
	}

	timer := metrics.NewTimer()
	metrics.Register("timer_name", timer)
	timer.Time(func() { time.Sleep(time.Millisecond * 10) })
	timer.Update(47)

	metrics.Write(metrics.DefaultRegistry, 1*time.Second, os.Stdout)
	//metrics.WriteJSONOnce(metrics.DefaultRegistry, os.Stdout)

	//report to open falcon
	//cfg := DefaultFalconConfig
	//cfg.Debug = true
	////report for each 5 seconds
	//cfg.Step = 5
	//falcon := NewFalcon(&cfg)
	//falcon.ReportRegistry(metrics.DefaultRegistry)

	//Metrics.Counter()
}
