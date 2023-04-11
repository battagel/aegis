package scanner

import (
	"aegis/pkg/logger"
)

type EndpointScanner struct {
	logger logger.Logger
}

func CreateEndpointScanner(logger logger.Logger) (*EndpointScanner, error) {
	logger.Debugln("Creating EndpointScanner")
	return &EndpointScanner{
		logger: logger,
	}, nil
}
