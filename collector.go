package main

import "sync"

// for collecting prometheus.constHistogram objects
type CustomCollector struct {
	//prometheus.Collector

	metric prometheus.Metric
	mutex  *sync.Mutex
}

func NewCustomCollector(mutex *sync.Mutex) *CustomCollector {
	return &CustomCollector{
		mutex: mutex,
	}
}

func (c *CustomCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	if c.metric != nil {
		val := c.metric
		ch <- val
	}
	c.mutex.Unlock()
}

func (p *CustomCollector) Describe(ch chan<- *prometheus.Desc) {
	// empty method to fulfill prometheus.Collector interface
}
