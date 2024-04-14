package saywo_metrics

import (
	"github.com/Edward-Alphonse/saywo_metrics/utils"
	"github.com/rcrowley/go-metrics"
)

func Start() {
	for _, reporter := range metric.Reporters {
		reporter.Start()
	}
}

func Close() {
	for _, reporter := range metric.Reporters {
		reporter.Close()
	}
}

func IncCounter(name string, labels map[string]string) {
	IncCounterWithStep(name, labels, 1)
}

func IncCounterWithStep(name string, labels map[string]string, step int64) {
	key := utils.EncodeKey(name, labels)
	counter := metrics.GetOrRegisterCounter(key, metrics.DefaultRegistry)
	counter.Inc(step)
}

func DecCounter(name string, labels map[string]string) {
	DecCounterWithStep(name, labels, 1)
}

func DecCounterWithStep(name string, labels map[string]string, step int64) {
	key := utils.EncodeKey(name, labels)
	counter := metrics.GetOrRegisterCounter(key, metrics.DefaultRegistry)
	counter.Dec(step)
}

func Gauge(name string, labels map[string]string, value int64) {
	key := utils.EncodeKey(name, labels)
	gauge := metrics.GetOrRegisterGauge(key, metrics.DefaultRegistry)
	gauge.Update(value)
}

func GaugeFloat64(name string, labels map[string]string, value float64) {
	key := utils.EncodeKey(name, labels)
	gauge := metrics.GetOrRegisterGaugeFloat64(key, metrics.DefaultRegistry)
	gauge.Update(value)
}

func Histogram(name string, labels map[string]string) {
	//key := utils.EncodeKey(name, labels)
	//metrics.GetOrRegisterHistogram(key, metrics.DefaultRegistry)
}

func Meter(name string, labels map[string]string) {
	key := utils.EncodeKey(name, labels)
	metrics.GetOrRegisterMeter(key, metrics.DefaultRegistry)
}

func Timer(name string, labels map[string]string) {
	key := utils.EncodeKey(name, labels)
	metrics.GetOrRegisterTimer(key, metrics.DefaultRegistry)
}
