package metrics

import (
	"aegis/internal/config"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	filesScanned  prometheus.Counter
	infectedFiles prometheus.Counter
	cleanFiles    prometheus.Counter
	scanErrors    prometheus.Counter
}

type Prometheus struct {
	metricChan chan string
	endpoint   string
	path       string
	metrics    Metrics
}

func CreateMetricManager(metricChan chan string) (*Prometheus, error) {
	fmt.Println("Creating Metric Server")
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config in metrics")
		return nil, err
	}
	endpoint := config.Services.Prometheus.Endpoint
	path := config.Services.Prometheus.Path

	filesScanned := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_total_scans", Help: ""})
	cleanFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_infected_files_scanned", Help: ""})
	infectedFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_clean_files_scanned", Help: ""})
	scanErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scan_errors", Help: ""})

	metrics := Metrics{filesScanned: filesScanned, infectedFiles: infectedFiles, cleanFiles: cleanFiles, scanErrors: scanErrors}

	return &Prometheus{metricChan: metricChan, endpoint: endpoint, path: path, metrics: metrics}, nil
}

func (p *Prometheus) StartMetricManager() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe("localhost:2112", nil)
	}()
	for {
		select {
		case metric := <-p.metricChan:
			switch metric {
			case "file_scanned":
				p.metrics.filesScanned.Inc()
			case "infected_file":
				p.metrics.infectedFiles.Inc()
			case "clean_file":
				p.metrics.cleanFiles.Inc()
			case "scan_error":
				p.metrics.scanErrors.Inc()
			}
		}
	}
}
