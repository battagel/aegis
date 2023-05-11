package events

import (
	"aegis/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type eventsCollector struct {
	logger           logger.Logger
	messagesReceived prometheus.Counter
	eventsError      prometheus.Counter
}

func CreateEventsCollector(logger logger.Logger) (*eventsCollector, error) {
	messagesReceieved := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_events_total_messages", Help: "Events total messages received"})
	eventsErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_events_total_errors", Help: "Events total errors"})
	return &eventsCollector{
		logger:           logger,
		messagesReceived: messagesReceieved,
		eventsError:      eventsErrors,
	}, nil
}

// Metric update functions
func (c *eventsCollector) MessageReceived() {
	c.logger.Debugln("Incrementing kafka message received counter")
	c.messagesReceived.Inc()
}

func (c *eventsCollector) EventsError() {
	c.logger.Debugln("Incrementing kafka event error counter")
	c.eventsError.Inc()
}
