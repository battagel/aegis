package dispatcher

import (
	"aegis/internal/object"
	"testing"
	"time"
)

type testScanner struct{}

func (s *testScanner) ScanObject(o *object.Object) error {
	return nil
}

func TestCreateDispatcher(t *testing.T) {
	scanner := &testScanner{}
	scanChan := make(chan *object.Object)
	dispatcher, err := CreateDispatcher([]Scanner{scanner}, scanChan)

	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}

	if dispatcher == nil {
		t.Error("expected dispatcher to be initialized, but got nil")
	}
}

func TestStartDispatcher(t *testing.T) {
	scanner := &testScanner{}
	scanChan := make(chan *object.Object)

	dispatcher, err := CreateDispatcher([]Scanner{scanner}, scanChan)
	if err != nil {
		t.Errorf("unexpected error while creating dispatcher: %v", err)
	}

	go dispatcher.StartDispatcher()

	objectToScan := &object.Object{}
	scanChan <- objectToScan

	// Sleep for some time to allow scanner to scan the object
	// and to prevent the test from exiting before the scanner completes
	// its task
	time.Sleep(time.Millisecond * 10)

	// Ensure that the ScanObject method of the scanner is called
	if objectToScan.ScannedBy != scanner {
		t.Error("expected object to be scanned by the scanner, but it wasn't")
	}

	// Stop the dispatcher loop
	dispatcher.StopDispatcher()
}

func TestStopDispatcher(t *testing.T) {
	scanner := &testScanner{}
	scanChan := make(chan *object.Object)
	dispatcher, err := CreateDispatcher([]Scanner{scanner}, scanChan)

	if err != nil {
		t.Errorf("unexpected error while creating dispatcher: %v", err)
	}

	err = dispatcher.StopDispatcher()

	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
}
