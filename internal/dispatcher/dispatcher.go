package dispatcher

import (
	"antivirus/internal/object"
	"fmt"
)

// TODO Double check this? I doubt embedding objectStore into object is correct
type Scanner interface {
	ScanObject(*object.Object) error
}

type Dispatcher struct {
	scanChan chan *object.Object
	scanners []Scanner
}

func CreateDispatcher(scanners []Scanner, scanChan chan *object.Object) (*Dispatcher, error) {
	fmt.Println("Creating Scan Manager")
	return &Dispatcher{scanChan: scanChan, scanners: scanners}, nil
}

func (d *Dispatcher) StartDispatcher() error {
	fmt.Println("Starting Scan Manager")
	for {
		request := <-d.scanChan
		// go s.scanObject(request.BucketName, request.ObjectKey)
		//
		for _, scanner := range d.scanners {
			// Ref?
			go scanner.ScanObject(request)
		}
	}
}

func (d *Dispatcher) StopDispatcher() error {
	fmt.Println("Stopping Scan Manager")

	return nil
}
