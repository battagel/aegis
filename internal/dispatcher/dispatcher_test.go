package dispatcher

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScanner is a mock implementation of the Scanner interface
type MockScanner struct {
	mock.Mock
}

// ScanObject is a mocked method for ScanObject in Scanner interface
func (m *MockScanner) ScanObject(obj *object.Object) error {
	args := m.Called(obj)
	return args.Error(0)
}

func TestDispatcher_Start(t *testing.T) {
	// Create a logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.NoError(t, err, "Error creating logger")

	// Create mock scanners
	mockScanner1 := new(MockScanner)
	mockScanner2 := new(MockScanner)

	// Create a dispatcher with mock scanners
	dispatcher, err := CreateDispatcher(logger, []Scanner{mockScanner1, mockScanner2}, make(chan *object.Object))
	assert.NoError(t, err, "Expected no error")

	// Expect that ScanObject will be called on both scanners with the same object
	obj := &object.Object{}
	mockScanner1.On("ScanObject", obj).Return(nil)
	mockScanner2.On("ScanObject", obj).Return(nil)

	// Start the dispatcher
	err = dispatcher.Start()
	assert.NoError(t, err, "Expected no error")

	// Assert that ScanObject was called on both scanners
	mockScanner1.AssertExpectations(t)
	mockScanner2.AssertExpectations(t)
}

func TestDispatcher_Stop(t *testing.T) {
	// Create a logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.NoError(t, err, "Error creating logger")

	// Create mock scanners
	mockScanner1 := new(MockScanner)
	mockScanner2 := new(MockScanner)

	// Create a dispatcher with mock scanners
	dispatcher, err := CreateDispatcher(logger, []Scanner{mockScanner1, mockScanner2}, make(chan *object.Object))
	assert.NoError(t, err, "Expected no error")

	// Call Stop on the dispatcher
	err = dispatcher.Stop()

	// Assert that Stop returned no error
	assert.NoError(t, err, "Expected no error")
}
