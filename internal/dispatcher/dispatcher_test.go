package dispatcher

import (
	"aegis/internal/object"
	"aegis/mocks"
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tableTest struct {
	bucketName, objectKey string
}

var tableTests = []tableTest{
	{"test", "test"},
	{"1234", "1234"},
}

func TestDispatcher_Start(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	mockScanner := new(mocks.Scanner)
	scanChan := make(chan *object.Object)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	dispatcher, err := CreateDispatcher(logger, []Scanner{mockScanner}, scanChan)
	assert.Nil(t, err)

	go dispatcher.Start(errChan, doneChan)

	for _, test := range tableTests {
		logger.Infow("Test started",
			"bucketName", test.bucketName,
			"objectKey", test.objectKey,
		)
		request, err := object.CreateObject(logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)
		mockScanner.On("ScanObject", request, errChan).Return(nil)
		scanChan <- request
		select {
		case err := <-errChan:
			logger.Errorw("Test failed",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
			t.Fail()
		default:
			logger.Infow("Test passed",
				"test", test,
			)
		}
	}
	logger.Infoln("Test Complete: Closing scanChan")
	close(scanChan)
}
