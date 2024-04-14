package main

import (
	"github.com/Edward-Alphonse/saywo_metrics"
	"github.com/Edward-Alphonse/saywo_metrics/sender"
	"os"
	"time"
)

func main() {
	config := &sender.ALiSLSConfig{
		DNS:             "cn-wuhan-lr.log.aliyuncs.com",
		AccessKeyId:     os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"),
		ProjectName:     "hzc-test",
		LogStoreName:    "metrics_test",
		Topic:           "test",
	}

	saywo_metrics.Register(saywo_metrics.ALiSLS(config))
	saywo_metrics.Start()

	saywo_metrics.Gauge("test", map[string]string{
		"123": "456",
	}, 2)

	saywo_metrics.Close()
	time.Sleep(10 * time.Second)
}
