package metrics

import (
	"aegis/pkg/logger"
)

type Prometheus interface {
	Start(chan error) error
	Stop() error
}

type MetricsManager struct {
	logger     logger.Logger
	prometheus Prometheus
}

func CreateMetricsManager(logger logger.Logger, prometheus Prometheus) (*MetricsManager, error) {
	logger.Debugln("Creating Metrics Manager")
	return &MetricsManager{
		logger:     logger,
		prometheus: prometheus,
	}, nil
}

func (m *MetricsManager) Start(errChan chan error) {
	m.logger.Debugln("Starting Metrics Manager")
	err := m.prometheus.Start()
	errChan <- err
}

func (m *MetricsManager) Stop() error {
	m.logger.Debugln("Stopping Metrics Manager")
	return m.prometheus.Stop()
}
