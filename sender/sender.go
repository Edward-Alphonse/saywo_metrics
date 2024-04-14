package sender

type MetricSender interface {
	Start()
	Close()
	SendMetric(name string, tags map[string]string, value float64, ts int64, source string) error
}
