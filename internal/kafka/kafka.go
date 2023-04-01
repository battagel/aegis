package kafka

import (
	"antivirus/internal/config"
	"antivirus/internal/object"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/segmentio/kafka-go"
)

type KafkaMgr struct {
	reader   *kafka.Reader
	scanChan chan *object.Object
}

func CreateKafkaManager(scanChan chan *object.Object) (*KafkaMgr, error) {
	fmt.Println("Creating Kafka Manager")
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config in kafka: ", err)
		return nil, err
	}
	conf := kafka.ReaderConfig{
		Brokers:  config.Services.Kafka.Brokers,
		Topic:    config.Services.Kafka.Topic,
		GroupID:  config.Services.Kafka.GroupID,
		MaxBytes: config.Services.Kafka.MaxBytes,
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
		newPut, bucketName, objectKey, err := k.decodeMessage(message)
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			return nil, err
		}
		if newPut {
			fmt.Println("Message: ", string(message.Value))
			request := object.CreateObject(bucketName, objectKey)
			k.scanChan <- request
		}
	}
}

func (k *KafkaMgr) StopKafkaManager() {
	fmt.Println("Stopping Kafka Consumer")
}

type MessageJson struct {
	EventName string `json:"EventName"`
	Records   []struct {
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

func (k *KafkaMgr) decodeMessage(message kafka.Message) (bool, string, string, error) {
	data := MessageJson{}
	err := json.Unmarshal([]byte(message.Value), &data)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
		return false, "", "", err
	}
	// Potential Security Issue. Attackers can edit tags without scan
	if data.EventName == "s3:ObjectCreated:PutTagging" {
		return false, "", "", nil
	}
	// Using url.QueryUnescape to handle spaces in object names as they show as "+" which breaks the file path
	bucketName, err := url.QueryUnescape(data.Records[0].S3.Bucket.Name)
	if err != nil {
		fmt.Println("Error unescaping bucket name: ", err)
		return false, "", "", err
	}
	objectKey, err := url.QueryUnescape(data.Records[0].S3.Object.Key)
	if err != nil {
		fmt.Println("Error unescaping object key: ", err)
		return false, "", "", err
	}
	return true, bucketName, objectKey, nil
}
