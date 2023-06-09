package metrics

import (
	"aegis/pkg/logger"
)

type MetricService interface {
	Start() error
	Stop() error
}

type MetricsManager struct {
	logger     logger.Logger
	prometheus MetricService
}

func CreateMetricsManager(logger logger.Logger, prometheus MetricService) (*MetricsManager, error) {
	logger.Debugln("Creating Metrics Manager")
	return &MetricsManager{
		logger:     logger,
		prometheus: prometheus,
	}, nil
}

func (m *MetricsManager) Start() error {
	m.logger.Debugln("Starting Metrics Manager")
	return m.prometheus.Start()
}

func (m *MetricsManager) Stop() error {
	m.logger.Debugln("Stopping Metrics Manager")
	return m.prometheus.Stop()
}
