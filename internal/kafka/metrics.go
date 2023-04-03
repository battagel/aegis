package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type kafkaCollector struct {
	sugar            *zap.SugaredLogger
	messagesReceived prometheus.Counter
}

func CreateKafkaCollector(sugar *zap.SugaredLogger) (*kafkaCollector, error) {
	messagesReceieved := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_kafka_total_messages", Help: "Kafka total messages received"})
	return &kafkaCollector{
		sugar:            sugar,
		messagesReceived: messagesReceieved,
	}, nil
}

// Metric update functions
func (c *kafkaCollector) MessageReceived() {
	c.sugar.Debug("Incrementing kafka message received counter")
	c.messagesReceived.Inc()
}
