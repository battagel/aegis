package events

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"context"
	// "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ctx = context.WithValue(context.Background(), "key", "value")

type tableTest struct {
	arg1     context.Context
	expected struct {
		bucketName string
		objectKey  string
	}
}

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

func TestEventsManager_Start_Happy(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	scanChan := make(chan *object.Object)
	errChan := make(chan error)
	mockKafka := new(MockEventsQueue)
	mockKafka.On("ReadMessage", ctx).Return("test", "test", nil)
	mockEventsCollector := new(MockEventsCollector)
	mockEventsCollector.On("MessageReceived").Return(nil)

	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
	assert.Nil(t, err)

	go eventsManager.Start(ctx, errChan)
	for _, test := range tableTests {
		request, err := object.CreateObject(logger, test.expected.bucketName, test.expected.objectKey)
		assert.Nil(t, err)
		assert.Equal(t, request, <-scanChan)
		select {
		case err := <-errChan:
			logger.Errorw("Test failed",
				"bucketName", test.expected.bucketName,
				"objectKey", test.expected.objectKey,
				"err", err,
			)
			t.Fail()
		default:
			logger.Infow("Test passed",
				"bucketName", test.expected.bucketName,
				"objectKey", test.expected.objectKey,
			)
		}
	}
}

func TestEventsManager_Close_Happy(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	closeCtx, closeCtxCancel := context.WithCancel(context.Background())
	scanChan := make(chan *object.Object)
	errChan := make(chan error)
	mockKafka := new(MockEventsQueue)
	mockKafka.On("ReadMessage", closeCtx).Return("test", "test", nil)
	mockEventsCollector := new(MockEventsCollector)
	mockEventsCollector.On("MessageReceived").Return(nil)

	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
	assert.Nil(t, err)

	go eventsManager.Start(closeCtx, errChan)
	time.Sleep(1 * time.Second)
	closeCtxCancel()
	time.Sleep(1 * time.Second)
	select {
	case err := <-errChan:
		logger.Errorw("Test failed",
			"err", err,
		)
		t.Fail()
	default:
		logger.Infoln("Test passed")
	}
}

// func TestEventsManager_Close_ErrorCloseKafka(t *testing.T) {
// 	// Create mock logger
// 	logger, err := logger.CreateZapLogger("debug", "console")
// 	assert.Nil(t, err)

// 	closeCtx, closeCtxCancel := context.WithCancel(context.Background())
// 	scanChan := make(chan *object.Object)
// 	errChan := make(chan error)
// 	mockKafka := new(MockEventsQueue)
// 	mockKafka.On("ReadMessage", closeCtx).Return("test", "test", nil)
// 	mockKafka.On("Close").Return(errors.New("Error Closing Kafka"))
// 	mockEventsCollector := new(MockEventsCollector)
// 	mockEventsCollector.On("MessageReceived").Return(nil)
// 	mockEventsCollector.On("EventsError").Return(nil)

// 	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
// 	assert.Nil(t, err)

// 	go eventsManager.Start(closeCtx, errChan)
// 	time.Sleep(1 * time.Second)
// 	logger.Infoln("Closing Kafka")
// 	closeCtxCancel()
// 	time.Sleep(1 * time.Second)
// 	select {
// 	case err := <-errChan:
// 		assert.EqualError(t, err, "Error Closing Kafka")
// 		logger.Infow("Test passed: Error recieved",
// 			"err", err,
// 		)
// 	default:
// 		logger.Errorln("Test failed")
// 		t.Fail()
// 	}
// }

// func TestEventsManager_Start_ErrorDecodeMsg(t *testing.T) {
// 	// Create mock logger
// 	logger, err := logger.CreateZapLogger("debug", "console")
// 	assert.Nil(t, err)

// 	scanChan := make(chan *object.Object)
// 	errChan := make(chan error)
// 	mockKafka := new(MockEventsQueue)
// 	mockKafka.On("ReadMessage", ctx).Return("", "", errors.New("Error Decoding Message"))
// 	mockEventsCollector := new(MockEventsCollector)
// 	mockEventsCollector.On("MessageReceived").Return(nil)
// 	mockEventsCollector.On("EventsError").Return(nil)

// 	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
// 	assert.Nil(t, err)

// 	go eventsManager.Start(ctx, errChan)
// 	for _, test := range tableTests {
// 		request, err := object.CreateObject(logger, test.expected.bucketName, test.expected.objectKey)
// 		assert.Nil(t, err)
// 		assert.Equal(t, request, <-scanChan)
// 		select {
// 		case err := <-errChan:
// 			assert.EqualError(t, err, "Error Decoding Message")
// 			logger.Errorw("Test passed: Error received",
// 				"bucketName", test.expected.bucketName,
// 				"objectKey", test.expected.objectKey,
// 				"err", err,
// 			)
// 		default:
// 			logger.Infow("Test failed",
// 				"bucketName", test.expected.bucketName,
// 				"objectKey", test.expected.objectKey,
// 			)
// 			t.Fail()
// 		}
// 	}
// }

// func TestEventsManager_Start_ErrorCreateObj(t *testing.T) {
// 	// Create mock logger
// 	logger, err := logger.CreateZapLogger("debug", "console")
// 	assert.Nil(t, err)

// 	scanChan := make(chan *object.Object)
// 	errChan := make(chan error)
// 	mockKafka := new(MockEventsQueue)
// 	mockKafka.On("ReadMessage", ctx).Return("", "", errors.New("Error Creating Object"))
// 	mockEventsCollector := new(MockEventsCollector)
// 	mockEventsCollector.On("MessageReceived").Return(nil)
// 	mockEventsCollector.On("EventsError").Return(nil)

// 	eventsManager, err := CreateEventsManager(logger, scanChan, mockKafka, mockEventsCollector)
// 	assert.Nil(t, err)

// 	go eventsManager.Start(ctx, errChan)
// 	for _, test := range tableTests {
// 		request, err := object.CreateObject(logger, test.expected.bucketName, test.expected.objectKey)
// 		assert.Nil(t, err)
// 		assert.Equal(t, request, <-scanChan)
// 		assert.Equal(t, request, <-scanChan)
// 		select {
// 		case err := <-errChan:
// 			assert.EqualError(t, err, "Error Creating Object")
// 			logger.Errorw("Test passed: Error received",
// 				"bucketName", test.expected.bucketName,
// 				"objectKey", test.expected.objectKey,
// 				"err", err,
// 			)
// 		default:
// 			logger.Infow("Test failed",
// 				"bucketName", test.expected.bucketName,
// 				"objectKey", test.expected.objectKey,
// 			)
// 			t.Fail()
// 		}
// 	}
// }
