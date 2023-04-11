package kafka

import (
	"aegis/pkg/logger"
	"context"
	"encoding/json"
	"net/url"

	"github.com/segmentio/kafka-go"
)

const (
	groupID  = "g1"
	maxBytes = 10
)

type KafkaConsumer struct {
	logger      logger.Logger
	kafkaReader *kafka.Reader
}

func CreateKafkaConsumer(logger logger.Logger, brokers []string, topic string) (*KafkaConsumer, error) {
	logger.Debugw("Creating Kafka Consumer",
		"brokers", brokers,
		"topic", topic,
	)
	readerConfig := kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: maxBytes,
	}
	kafkaReader := kafka.NewReader(readerConfig)
	return &KafkaConsumer{
		logger:      logger,
		kafkaReader: kafkaReader,
	}, nil
}

func (k *KafkaConsumer) ReadMessage() (string, string, error) {
	k.logger.Debugw("Reading message from Kafka")
	message, err := k.kafkaReader.ReadMessage(context.Background())
	if err != nil {
		k.logger.Errorw("Error reading message from Kafka",
			"error", err,
		)
		return "", "", err
	}
	bucketName, objectKey, err := k.decodeMessage(message)
	if err != nil {
		k.logger.Errorw("Error decoding message",
			"error", err,
		)
		return "", "", err
	}
	return bucketName, objectKey, nil
}

func (k *KafkaConsumer) Close() error {
	k.logger.Debugw("Closing Kafka Consumer")
	err := k.kafkaReader.Close()
	return err
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

func (k *KafkaConsumer) decodeMessage(message kafka.Message) (string, string, error) {
	data := MessageJson{}
	err := json.Unmarshal([]byte(message.Value), &data)
	if err != nil {
		k.logger.Errorw("Error unmarshalling json",
			"error", err,
		)
		return "", "", err
	}
	// Potential Security Issue. Attackers can edit tags without scan
	if data.EventName == "s3:ObjectCreated:PutTagging" {
		return "", "", nil
	}
	k.logger.Infow("Message",
		"message", string(message.Value),
	)
	// Using url.QueryUnescape to handle spaces in object names as they show as "+" which breaks the file path
	bucketName, err := url.QueryUnescape(data.Records[0].S3.Bucket.Name)
	if err != nil {
		k.logger.Errorw("Error unescaping bucket name",
			"error", err,
		)
		return "", "", err
	}
	objectKey, err := url.QueryUnescape(data.Records[0].S3.Object.Key)
	if err != nil {
		k.logger.Errorw("Error unescaping object key",
			"error", err,
		)
		return "", "", err
	}
	return bucketName, objectKey, nil
}
