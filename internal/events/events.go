package events

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
)

type EventsCollector interface {
	MessageReceived()
}

type Kafka interface {
	ReadMessage() (string, string, error)
	Close() error
}

type EventsManager struct {
	logger          logger.Logger
	kafka           Kafka
	scanChan        chan *object.Object
	eventsCollector EventsCollector
}

func CreateEventsManager(logger logger.Logger, scanChan chan *object.Object, kafka Kafka, eventsCollector EventsCollector) (*EventsManager, error) {
	logger.Debugln("Creating Event Manager")
	return &EventsManager{
		logger:          logger,
		kafka:           kafka,
		scanChan:        scanChan,
		eventsCollector: eventsCollector,
	}, nil
}

func (k *EventsManager) Start(errChan chan error) (*EventsManager, error) {
	k.logger.Debugln("Listening for activity on Kafka...")
	for {
		bucketName, objectKey, err := k.kafka.ReadMessage()
		if err != nil {
			k.logger.Errorw("Error decoding message",
				"error", err,
			)
			errChan <- err
			return nil, err
		}
		if bucketName != "" && objectKey != "" {
			k.eventsCollector.MessageReceived()
			request, err := object.CreateObject(k.logger, bucketName, objectKey)
			if err != nil {
				k.logger.Errorw("Error creating object",
					"error", err,
				)
				errChan <- err
				return nil, err
			}
			k.scanChan <- request
		}
	}
}

func (k *EventsManager) Stop() error {
	k.logger.Debugln("Stopping Kafka Consumer")
	err := k.kafka.Close()
	return err
}
