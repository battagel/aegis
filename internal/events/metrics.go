package events

import (
	"aegis/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type eventsCollector struct {
	logger           logger.Logger
	messagesReceived prometheus.Counter
}

func CreateEventsCollector(logger logger.Logger) (*eventsCollector, error) {
	messagesReceieved := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_kafka_total_messages", Help: "Kafka total messages received"})
	return &eventsCollector{
		logger:           logger,
		messagesReceived: messagesReceieved,
	}, nil
}

// Metric update functions
func (c *eventsCollector) MessageReceived() {
	c.logger.Debugln("Incrementing kafka message received counter")
	c.messagesReceived.Inc()
}
