package dispatcher

import (
	"aegis/internal/object"
	"fmt"
)

type Scanner interface {
	ScanObject(*object.Object) error
}

type Dispatcher struct {
	scanChan chan *object.Object
	scanners []Scanner
}

func CreateDispatcher(scanners []Scanner, scanChan chan *object.Object) (*Dispatcher, error) {
	fmt.Println("Creating dispatcher")
	return &Dispatcher{scanChan: scanChan, scanners: scanners}, nil
}

func (d *Dispatcher) StartDispatcher() error {
	fmt.Println("Starting dispatcher loop")
	for {
		request := <-d.scanChan
		for _, scanner := range d.scanners {
			// Ref?
			go scanner.ScanObject(request)
		}
	}
}

func (d *Dispatcher) StopDispatcher() error {
	fmt.Println("Stopping dispatcher")
	return nil
}
