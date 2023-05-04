package events

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"context"
)

type EventsCollector interface {
	MessageReceived()
}

type Kafka interface {
	ReadMessage(context.Context) (string, string, error)
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

func (k *EventsManager) Start(ctx context.Context, errChan chan error) {
	k.logger.Debugln("Listening for activity on Kafka...")
	for {
		select {
		case <-ctx.Done():
			k.logger.Debugln("Stopping Kafka consumer")
			close(k.scanChan)
			err := k.kafka.Close()
			if err != nil {
				k.logger.Errorw("Error closing kafka consumer",
					"error", err,
				)
				errChan <- err
			}
			return
		default:
			bucketName, objectKey, err := k.kafka.ReadMessage(ctx)
			if err != nil {
				k.logger.Errorw("Error decoding message",
					"error", err,
				)
				errChan <- err
				continue
			}
			if bucketName != "" && objectKey != "" {
				k.eventsCollector.MessageReceived()
				request, err := object.CreateObject(k.logger, bucketName, objectKey)
				if err != nil {
					k.logger.Errorw("Error creating object",
						"error", err,
					)
					errChan <- err
					continue
				}
				k.scanChan <- request
			}
		}
	}
}
