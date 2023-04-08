package dispatcher

import (
	"aegis/internal/object"
	"aegis/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type CommonTestItems struct {
	Sugar *zap.SugaredLogger
}

func ProvideCommonTestItems(t *testing.T) *CommonTestItems {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Error creating logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return &CommonTestItems{Sugar: sugar}
}

func TestDispatcher(t *testing.T) {
	commonTestItems := ProvideCommonTestItems(t)
	scanChan := make(chan *object.Object)
	mockScanner := new(mocks.Scanner)

	testObjectGood := object.Object{
		ObjectKey:  "good.txt",
		BucketName: "test-bucket",
		CachePath:  "/cache",
	}

	testObjectBad := object.Object{
		ObjectKey:  "bad.txt",
		BucketName: "test-bucket",
		CachePath:  "/cache",
	}

	mockScanner.On("ScanObject", testObjectGood).Return(false)
	mockScanner.On("ScanObject", testObjectBad).Return(true)

	dispatcher, err := CreateDispatcher(commonTestItems.Sugar, []Scanner{mockScanner}, scanChan)
	assert.NoError(t, err)

	go dispatcher.StartDispatcher()

	// Send stuff to get scanned
	scanChan <- &testObjectGood
	scanChan <- &testObjectBad
	dispatcher.StopDispatcher()
}
