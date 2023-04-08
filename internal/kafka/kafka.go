package kafka

import (
	"aegis/internal/config"
	"aegis/internal/object"
	"context"
	"encoding/json"
	"net/url"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaCollector interface {
	MessageReceived()
}

type KafkaReader interface {
	ReadMessage(context.Context) (kafka.Message, error)
}

type KafkaMgr struct {
	sugar          *zap.SugaredLogger
	kafkaReader    KafkaReader
	scanChan       chan *object.Object
	kafkaCollector KafkaCollector
}

func CreateKafkaManager(sugar *zap.SugaredLogger, scanChan chan *object.Object, kafkaReader KafkaReader, kafkaCollector KafkaCollector) (*KafkaMgr, error) {
	config, err := config.GetConfig()
	if err != nil {
		sugar.Errorw("Error getting config in kafka",
			"error", err,
		)
		return nil, err
	}
	sugar.Debugw("Creating Kafka Manager",
		"brokers", config.Services.Kafka.Brokers,
		"topic", config.Services.Kafka.Topic,
	)
	return &KafkaMgr{
		sugar:          sugar,
		kafkaReader:    kafkaReader,
		scanChan:       scanChan,
		kafkaCollector: kafkaCollector,
	}, nil
}

func (k *KafkaMgr) StartKafkaManager() (*KafkaMgr, error) {
	k.sugar.Debugln("Listening for activity on Kafka...")
	for {
		message, err := k.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			k.sugar.Errorw("Error reading message from Kafka",
				"error", err,
			)
			return nil, err
		}
		newPut, bucketName, objectKey, err := k.decodeMessage(message)
		if err != nil {
			k.sugar.Errorw("Error decoding message",
				"error", err,
			)
			return nil, err
		}
		if newPut {
			k.sugar.Infow("Message",
				"message", string(message.Value),
			)
			k.kafkaCollector.MessageReceived()
			request, err := object.CreateObject(k.sugar, bucketName, objectKey)
			if err != nil {
				k.sugar.Errorw("Error creating object",
					"error", err,
				)
				return nil, err
			}
			k.scanChan <- request
		}
	}
}

func (k *KafkaMgr) StopKafkaManager() {
	k.sugar.Debugln("Stopping Kafka Consumer")
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
		k.sugar.Errorw("Error unmarshalling json",
			"error", err,
		)
		return false, "", "", err
	}
	// Potential Security Issue. Attackers can edit tags without scan
	if data.EventName == "s3:ObjectCreated:PutTagging" {
		return false, "", "", nil
	}
	// Using url.QueryUnescape to handle spaces in object names as they show as "+" which breaks the file path
	bucketName, err := url.QueryUnescape(data.Records[0].S3.Bucket.Name)
	if err != nil {
		k.sugar.Errorw("Error unescaping bucket name",
			"error", err,
		)
		return false, "", "", err
	}
	objectKey, err := url.QueryUnescape(data.Records[0].S3.Object.Key)
	if err != nil {
		k.sugar.Errorw("Error unescaping object key",
			"error", err,
		)
		return false, "", "", err
	}
	return true, bucketName, objectKey, nil
}
