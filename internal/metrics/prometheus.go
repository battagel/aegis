package metrics

import (
	"aegis/internal/config"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Prometheus struct {
	sugar    *zap.SugaredLogger
	endpoint string
	path     string
}

func CreateMetricManager(sugar *zap.SugaredLogger) (*Prometheus, error) {
	sugar.Debugln("Creating Metric Server")
	config, err := config.GetConfig()
	if err != nil {
		sugar.Errorw("Error getting config in metrics",
			"error", err,
		)
		return nil, err
	}
	endpoint := config.Services.Prometheus.Endpoint
	path := config.Services.Prometheus.Path

	return &Prometheus{sugar: sugar, endpoint: endpoint, path: path}, nil
}

func (p *Prometheus) StartMetricManager() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe("localhost:2112", nil)
	}()
}
