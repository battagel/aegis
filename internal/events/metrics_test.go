package events

import (
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanCollector_Happy(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	eventsCollector, err := CreateEventsCollector(logger)
	assert.Nil(t, err)

	eventsCollector.MessageReceived()
	eventsCollector.EventsError()
}
