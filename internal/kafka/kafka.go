package kafka

import (
	"antivirus/internal/object"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaMgr struct {
	reader   *kafka.Reader
	scanChan chan *object.Object
}

func CreateKafkaManager(topic string, scanChan chan *object.Object) (*KafkaMgr, error) {
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

func (k *KafkaMgr) StartKafkaManager() (*KafkaMgr, error) {
	fmt.Println("Listening for activity on Kafka...")
	for {
		message, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message from Kafka: ", err)
			return nil, err
		}
		fmt.Println("Message: ", string(message.Value))
		bucketName, objectKey, err := k.decodeMessage(message)
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			return nil, err
		}
		request := object.CreateObject(bucketName, objectKey)
		k.scanChan <- request
	}
}

func (k *KafkaMgr) StopKafkaManager() {
	fmt.Println("Stopping Kafka Consumer")
}

type MessageJson struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func (k *KafkaMgr) decodeMessage(message kafka.Message) (string, string, error) {
	data := MessageJson{}
	err := json.Unmarshal([]byte(message.Value), &data)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
		return "", "", err
	}
	bucketName := data.Records[0].S3.Bucket.Name
	objectKey := data.Records[0].S3.Object.Key
	// TODO Fix this json shit
	//objectPath := "test-bucket/42154d0805933548da9b7a9fbbce40be9e155091e6f96ed4ce324c21b3430b20"
	// objectPath := "test-bucket/Wave Ticket.pdf"

	return bucketName, objectKey, nil
}
