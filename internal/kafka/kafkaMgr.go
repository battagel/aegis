package kafkaMgr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaMgr struct {
	reader   *kafka.Reader
	scanChan chan string
}

func CreateKafkaMgr(topic string, scanChan chan string) (*KafkaMgr, error) {
	fmt.Println("Creating Kafka Manager")
	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  "g1",
		MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)
	return &KafkaMgr{reader: reader, scanChan: scanChan}, nil
}

func (k *KafkaMgr) StartKafkaMgr() (*KafkaMgr, error) {
	fmt.Println("Listening for activity on Kafka...")
	for {
		message, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message from Kafka: ", err)
			return nil, err
		}
		fmt.Println("Message: ", string(message.Value))
		objectPath, err := k.decodeMessage(message)
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			return nil, err
		}
		k.scanChan <- objectPath
	}
}

func (k *KafkaMgr) StopKafkaMgr() {
	fmt.Println("Stopping Kafka Consumer")
}

func (k *KafkaMgr) decodeMessage(message kafka.Message) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message.Value), &data)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
		return "", err
	}
	// TODO Fix this json shit
	objectPath := "test-bucket/Wave Ticket.pdf"

	return objectPath, nil
}
