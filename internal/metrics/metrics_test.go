package metrics

import (
	"aegis/mocks"
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsManager_Happy(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")

	mockPrometheus := new(mocks.MetricService)
	mockPrometheus.On("Start").Return(nil)
	mockPrometheus.On("Stop").Return(nil)

	metrics, err := CreateMetricsManager(logger, mockPrometheus)
	assert.Nil(t, err)

	err = metrics.Start()
	assert.Nil(t, err)

	err = metrics.Stop()
	assert.Nil(t, err)
}
