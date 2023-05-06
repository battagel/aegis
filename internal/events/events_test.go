package events

import (
	"aegis/internal/object"
	"aegis/mocks"
	"aegis/pkg/logger"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafka is a mock implementation of the Kafka interface
type MockKafka struct {
	mock.Mock
}

// MockEventsCollector is a mock implementation of the EventsCollector interface
type MockEventsCollector struct {
	mock.Mock
}

type tableTest struct {
	arg1     context.Context
	expected struct {
		bucketName string
		objectKey  string
	}
}

var ctx = context.WithValue(context.Background(), "key", "value")

var tableTests = []tableTest{
	{
		arg1: ctx,
		expected: struct {
			bucketName string
			objectKey  string
		}{
			bucketName: "test",
			objectKey:  "test",
		},
	},
}

func TestEventsManager_Start(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	scanChan := make(chan *object.Object)
	errChan := make(chan error)
	mockKafka := new(mocks.Kafka)
	mockKafka.On("ReadMessage", ctx).Return("test", "test", nil)
	mockEventsCollector := new(mocks.EventsCollector)
	mockEventsCollector.On("MessageReceived").Return(nil)

	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
	assert.Nil(t, err)

	go eventsManager.Start(ctx, errChan)
	for _, test := range tableTests {
		request, err := object.CreateObject(logger, test.expected.bucketName, test.expected.objectKey)
		assert.Nil(t, err)
		assert.Equal(t, request, <-scanChan)
	}
}

// func TestEventsManager_Stop(t *testing.T) {
// 	logger.CreateZapLogger("debug", "console")
// }
