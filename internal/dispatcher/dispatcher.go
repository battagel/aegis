package dispatcher

import (
	"aegis/internal/object"

	"go.uber.org/zap"
)

type Scanner interface {
	ScanObject(*object.Object) error
}

type Dispatcher struct {
	sugar    *zap.SugaredLogger
	scanChan chan *object.Object
	scanners []Scanner
}

func CreateDispatcher(sugar *zap.SugaredLogger, scanners []Scanner, scanChan chan *object.Object) (*Dispatcher, error) {
	sugar.Debugln("Creating dispatcher")
	return &Dispatcher{
		sugar:    sugar,
		scanChan: scanChan,
		scanners: scanners,
	}, nil
}

func (d *Dispatcher) StartDispatcher() error {
	d.sugar.Debugln("Starting dispatcher loop...")
	for {
		request := <-d.scanChan
		for _, scanner := range d.scanners {
			// Pass by reference??
			go scanner.ScanObject(request)
		}
	}
}

func (d *Dispatcher) StopDispatcher() error {
	d.sugar.Debugln("Stopping dispatcher")
	return nil
}
