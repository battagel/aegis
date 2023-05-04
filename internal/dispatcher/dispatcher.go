package dispatcher

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"sync"
)

type Scanner interface {
	ScanObject(*object.Object, chan error)
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

func (d *Dispatcher) Start(errChan chan error, done chan struct{}) {
	d.logger.Debugln("Starting dispatcher loop...")
	var wg sync.WaitGroup

	for request := range d.scanChan {
		for _, scanner := range d.scanners {
			wg.Add(1)
			go func(req *object.Object, sc Scanner) {
				defer wg.Done()
				sc.ScanObject(req, errChan)
			}(request, scanner)
		}
	}
	wg.Wait()
	d.logger.Debugln("Dispatcher loop stopped")
	done <- struct{}{}
}
