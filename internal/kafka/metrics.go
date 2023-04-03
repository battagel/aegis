package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type kafkaCollector struct {
	messagesReceived prometheus.Counter
}

func CreateKafkaCollector() (*kafkaCollector, error) {
	messagesReceieved := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_kafka_total_messages", Help: "Kafka total messages received"})
	return &kafkaCollector{
		messagesReceived: messagesReceieved,
	}, nil
}

// Metric update functions
func (c *kafkaCollector) MessageReceived() {
	c.messagesReceived.Inc()
}
