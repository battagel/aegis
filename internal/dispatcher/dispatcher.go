package dispatcher

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
)

type Scanner interface {
	ScanObject(*object.Object) error
}

type Dispatcher struct {
	logger   logger.Logger
	scanChan chan *object.Object
	scanners []Scanner
}

func CreateDispatcher(logger logger.Logger, scanners []Scanner, scanChan chan *object.Object) (*Dispatcher, error) {
	logger.Debugln("Creating dispatcher")
	return &Dispatcher{
		logger:   logger,
		scanChan: scanChan,
		scanners: scanners,
	}, nil
}

func (d *Dispatcher) Start() error {
	d.logger.Debugln("Starting dispatcher loop...")
	for {
		request := <-d.scanChan
		for _, scanner := range d.scanners {
			// Pass by reference??
			go scanner.ScanObject(request)
		}
	}
}

func (d *Dispatcher) Stop() error {
	d.logger.Debugln("Stopping dispatcher")
	return nil
}
