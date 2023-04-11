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

func CreatePrometheusServer(logger logger.Logger, endpoint, path string) (*Prometheus, error) {
	logger.Debugln("Creating Metric Server")
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())
	// TODO Add https support
	return &Prometheus{
		logger: logger,
		httpServer: &http.Server{
			Addr:         endpoint,
			ReadTimeout:  readTimeoutInSeconds * time.Second,
			WriteTimeout: writeTimeoutInSeconds * time.Second,
			Handler:      mux,
		},
	}, nil
}

func (p *Prometheus) Start() error {
	p.logger.Debugw("Exposing metrics at",
		"endpoint", p.httpServer.Addr,
		// TODO fix this error. Not finding path
		"path", p.httpServer.Handler,
	)

	// TODO: Add graceful shutdown this returns error when stopped
	err := p.httpServer.ListenAndServe()
	if err != nil {
		p.logger.Errorw("Error starting prometheus server",
			"error", err,
		)
		return err
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
