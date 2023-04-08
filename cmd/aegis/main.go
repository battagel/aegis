package main

import (
	"aegis/internal/cli"
	"aegis/internal/config"
	"aegis/internal/dispatcher"
	"aegis/internal/kafka"
	"aegis/internal/metrics"
	"aegis/internal/object"
	"aegis/internal/objectstore"
	"aegis/internal/scanner"
	"fmt"
	"os"
	"os/signal"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	kafkaGo "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func run() int {
	// Config and Logger
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config in main :", err)
	}
	var logger *zap.Logger
	if config.Debug {
		cli.PrintSplash()
		logger, err = zap.NewDevelopment()
		if err != nil {
			fmt.Println("Error creating logger in main :", err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			fmt.Println("Error creating logger in main :", err)
		}
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Removes hidden control flow
	scanChan := make(chan *object.Object)
	defer close(scanChan)

	sugar.Infoln("Starting Aegis")
	metricManager, err := metrics.CreateMetricManager(sugar)
	if err != nil {
		sugar.Errorw("Error creating metric server",
			"error", err,
		)
	}

	endpoint := config.Services.Minio.Endpoint
	accessKey := config.Services.Minio.AccessKey
	secretKey := config.Services.Minio.SecretKey
	useSSL := config.Services.Minio.UseSSL
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		sugar.Errorw("Connecting to MinIO failed",
			"error", err,
		)
	}
	objectStoreCollector, err := objectstore.CreateObjectStoreCollector(sugar)
	if err != nil {
		sugar.Errorw("Error creating object store collector",
			"error", err,
		)
	}
	objectStore, err := objectstore.CreateObjectStore(sugar, minioClient, objectStoreCollector)
	if err != nil {
		sugar.Errorw("Error creating object store",
			"error", err,
		)
	}

	conf := kafkaGo.ReaderConfig{
		Brokers:  config.Services.Kafka.Brokers,
		Topic:    config.Services.Kafka.Topic,
		GroupID:  config.Services.Kafka.GroupID,
		MaxBytes: config.Services.Kafka.MaxBytes,
	}
	kafkaReader := kafkaGo.NewReader(conf)

	kafkaCollector, err := kafka.CreateKafkaCollector(sugar)
	if err != nil {
		sugar.Errorw("Error creating kafka collector",
			"error", err,
		)
	}
	kafkaManager, err := kafka.CreateKafkaManager(sugar, scanChan, kafkaReader, kafkaCollector)
	if err != nil {
		sugar.Errorw("Error creating kafka manager",
			"error", err,
		)
	}

	scanCollector, err := scanner.CreateScanCollector(sugar)
	if err != nil {
		sugar.Errorw("Error creating scan collector",
			"error", err,
		)
	}
	clamAV, err := scanner.CreateClamAV(sugar, objectStore, scanCollector)
	if err != nil {
		sugar.Errorw("Error creating clamAV scanner",
			"error", err,
		)
	}

	dispatcher, err := dispatcher.CreateDispatcher(sugar, []dispatcher.Scanner{clamAV}, scanChan)
	if err != nil {
		sugar.Errorw("Error creating antivirus manager",
			"error", err,
		)
	}

	// sync.WaitGroup() as part of termination
	go kafkaManager.StartKafkaManager()
	go dispatcher.StartDispatcher()
	go metricManager.StartMetricManager()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan // Wait until interrupt
	sugar.Infoln("Shutting down Aegis")
	// Cleanup stuff ...
	kafkaManager.StopKafkaManager()
	dispatcher.StopDispatcher()
	// Only stop when all scans finished?
	// Send signals to kafka and scan maanger to stop
	// sync.waitgroup to wait for all scans to finish
	return 0
}

func main() {
	os.Exit(run())
}
