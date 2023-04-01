package metrics

import (
	"antivirus/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	metricChan chan string
	endpoint   string
	path       string
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
	return &Prometheus{metricChan: metricChan, endpoint: endpoint, path: path}, nil
}

var opsProcessed = promauto.NewCounter(prometheus.CounterOpts{Name: "antivirus_processed_ops_total", Help: "The total number of processed events"})

func (m *Prometheus) StartMetricManager() {
	fmt.Println("Starting Metric Server")

	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("localhost:2112", nil)
}
