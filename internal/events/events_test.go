package events

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafka is a mock implementation of the Kafka interface
type MockKafka struct {
	mock.Mock
}

// ReadMessage is a mock method for reading a Kafka message
func (m *MockKafka) ReadMessage() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

// Close is a mock method for closing the Kafka consumer
func (m *MockKafka) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockEventsCollector is a mock implementation of the EventsCollector interface
type MockEventsCollector struct {
	mock.Mock
}

// MessageReceived is a mock method for simulating a received message event
func (m *MockEventsCollector) MessageReceived() {
	m.Called()
}

func TestEventsManager_Start(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.NoError(t, err, "Error creating logger")

	// Create mock Kafka
	mockKafka := new(MockKafka)

	// Create mock scan channel
	scanChan := make(chan *object.Object)

	// Create mock events collector
	mockEventsCollector := new(MockEventsCollector)

	// Create EventsManager with mocks
	eventsManager := &EventsManager{
		logger:          logger,
		kafka:           mockKafka,
		scanChan:        scanChan,
		eventsCollector: mockEventsCollector,
	}

	// Set expectations for Kafka ReadMessage
	mockKafka.On("ReadMessage").Return("bucketName", "objectKey", nil).Once()

	// Set expectations for EventsCollector MessageReceived
	mockEventsCollector.On("MessageReceived").Once()

	// Invoke Start() method on EventsManager
	_, err = eventsManager.Start()

	// Assert that the expected methods were called and no error occurred
	mockKafka.AssertExpectations(t)
	mockEventsCollector.AssertExpectations(t)
	assert.NoError(t, err, "Expected no error")
}

func TestEventsManager_Stop(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.NoError(t, err, "Error creating logger")

	// Create mock Kafka
	mockKafka := new(MockKafka)

	// Create mock scan channel
	scanChan := make(chan *object.Object)

	// Create mock events collector
	mockEventsCollector := new(MockEventsCollector)

	// Create EventsManager with mocks
	eventsManager := &EventsManager{
		logger:          logger,
		kafka:           mockKafka,
		scanChan:        scanChan,
		eventsCollector: mockEventsCollector,
	}

	// Set expectations for Kafka Close
	mockKafka.On("Close").Return(nil).Once()

	// Invoke Stop() method on EventsManager
	err = eventsManager.Stop()

	// Assert that the expected methods were called and no error occurred
	mockKafka.AssertExpectations(t)
	assert.NoError(t, err, "Expected no error")
}
