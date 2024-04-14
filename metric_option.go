package saywo_metrics

import (
	"github.com/Edward-Alphonse/saywo_metrics/reporter"
	"github.com/Edward-Alphonse/saywo_metrics/sender"
	"time"
)

var metric *SwMetric

type SwMetric struct {
	Reporters []reporter.Reporter
}

type MetricOption interface {
	registryTo(metric *SwMetric)
}

type MetricOptionFunc func() reporter.Reporter

func (f MetricOptionFunc) registryTo(metric *SwMetric) {
	reporters := metric.Reporters
	if len(reporters) == 0 {
		reporters = make([]reporter.Reporter, 0)
	}
	metric.Reporters = append(reporters, f())
}

func ALiSLS(config *sender.ALiSLSConfig) MetricOption {
	slsSender := sender.NewALiSLSSender(config)
	return MetricOptionFunc(func() reporter.Reporter {
		setter := make([]reporter.Option, 0)
		if config.SampleInterval > 0 {
			interval := time.Duration(config.SampleInterval) * time.Second
			setter = append(setter, reporter.Interval(interval))
		}
		return reporter.NewMetricsReporter(slsSender, setter...)
	})
}

func Register(opts ...MetricOption) {
	metric = new(SwMetric)
	for _, opt := range opts {
		opt.registryTo(metric)
	}
}
