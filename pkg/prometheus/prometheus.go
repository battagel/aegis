package prometheus

import (
	"aegis/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	readTimeoutInSeconds   = 5
	writeTimeoutInSeconds  = 10
	serverTimeoutInSeconds = 5
)

type Prometheus struct {
	logger     logger.Logger
	httpServer *http.Server
}

func CreatePrometheusExporter(logger logger.Logger, endpoint, path string) (*Prometheus, error) {
	logger.Debugln("Creating Prometheus Server")
	return &Prometheus{
		logger: logger,
		httpServer: &http.Server{
			Addr:         endpoint,
			ReadTimeout:  readTimeoutInSeconds * time.Second,
			WriteTimeout: writeTimeoutInSeconds * time.Second,
			Handler:      promhttp.Handler(),
		},
	}, nil
}

func (p *Prometheus) Start() error {
	p.logger.Debugw("Exposing metrics at",
		"endpoint", p.httpServer.Addr,
	)

	err := p.httpServer.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			p.logger.Debugln("Prometheus server closed, exiting")
		} else {
			p.logger.Errorw("Error starting prometheus server",
				"error", err,
			)
			return err
		}
	}
	return nil
}

func (p *Prometheus) Stop() error {
	p.logger.Debugln("Stopping Metric Server")
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutInSeconds*time.Second)
	defer cancel()

	if err := p.httpServer.Shutdown(ctx); err != nil {
		p.logger.Errorw("Error stopping prometheus server",
			"error", err,
		)
		return err
	}
	return nil
}
