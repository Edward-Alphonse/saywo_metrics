package sender

import (
	"encoding/json"
	"log"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) Start() {}

func (s *ConsoleSender) Close() {}

func (s *ConsoleSender) SendMetric(name string, tags map[string]string, value float64, ts int64, source string) error {
	bytes, _ := json.Marshal(tags)

	log.Printf("metrics, name: %s, tags: %v, value: %f, timestamp: %d, source: %s", name, string(bytes), value, ts, source)
	return nil
}
