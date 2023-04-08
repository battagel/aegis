package kafka

import (
	"aegis/internal/object"
	"aegis/mocks"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	kafkaGo "github.com/segmentio/kafka-go"
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

func TestKafkaManager(t *testing.T) {
	commonTestItems := ProvideCommonTestItems(t)
	scanChan := make(chan *object.Object)
	mockKafkaCollector := new(mocks.KafkaCollector)
	mockKafkaReader := new(mocks.KafkaReader)

	kafka, err := CreateKafkaManager(commonTestItems.Sugar, scanChan, mockKafkaReader, mockKafkaCollector)
	assert.NoError(t, err)

	jsonFile, err := os.Open("examplemessage.json")
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	assert.NoError(t, err)

	exampleMessage := kafkaGo.Message{Value: jsonBytes}
	mockKafkaReader.On("ReadMessage").Return(exampleMessage, nil)

	go kafka.StartKafkaManager()

	kafka.StopKafkaManager()
}
