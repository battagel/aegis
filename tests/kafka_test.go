package kafka

import (
	"antivirus/internal/kafka"
	"testing"
)

func TestKafka(t *testing.T) {
	fmt.Println("Testing Kafka Manager")

	kafka, err := kafka.StartKafkaConsumer()

	if err != nil {
		t.Error(err)
	}
}
