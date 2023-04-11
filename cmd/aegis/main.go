package main

import (
	"aegis/internal/auditlog"
	"aegis/internal/cli"
	"aegis/internal/dispatcher"
	"aegis/internal/events"
	"aegis/internal/object"
	"aegis/internal/objectstore"
	"aegis/internal/scanner"
	"aegis/pkg/clamav"
	"aegis/pkg/config"
	"aegis/pkg/kafka"
	"aegis/pkg/logger"
	"aegis/pkg/minio"
	"aegis/pkg/postgres"
	"aegis/pkg/prometheus"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	cli.PrintSplash()
	config, err := config.GetConfig()
	logger, err := logger.CreateZapLogger(config.Logger.Level, config.Logger.Encoding)
	if err != nil {
		fmt.Println("Error creating logger", err)
	}

	scanChan := make(chan *object.Object)
	defer close(scanChan)

	logger.Infoln("Starting Aegis")
	logger.Debugln("Creating metric collectors")
	metrics, err := prometheus.CreatePrometheusServer(logger, config.Prometheus.Endpoint, config.Prometheus.Path)
	if err != nil {
		logger.Errorw("Error creating metric collectors",
			"error", err,
		)
	}
	objectStoreCollector, err := objectstore.CreateObjectStoreCollector(logger)
	if err != nil {
		logger.Errorw("Error creating object store collector",
			"error", err,
		)
	}
	eventsCollector, err := events.CreateKafkaCollector(logger)
	if err != nil {
		logger.Errorw("Error creating kafka collector",
			"error", err,
		)
	}
	scanCollector, err := scanner.CreateScanCollector(logger)
	if err != nil {
		logger.Errorw("Error creating collectors",
			"error", err,
		)
	}

	postgresDB, dbClose, err := postgres.CreatePostgresDB(logger, config.Postgres.User, config.Postgres.Password, config.Postgres.Endpoint, config.Postgres.Database)
	if err != nil {
		logger.Errorw("Error creating postgres database",
			"error", err,
		)
	}
	defer dbClose()
	auditLogger, err := auditlog.CreateAuditLogger(logger, postgresDB, config.Postgres.Table)

	minioStore, err := minio.CreateMinio(logger, config.Minio.Endpoint, config.Minio.AccessKey, config.Minio.SecretKey, config.Minio.UseSSL)
	if err != nil {
		logger.Errorw("Error creating minio client",
			"error", err,
		)
	}
	objectStore, err := objectstore.CreateObjectStore(logger, minioStore, objectStoreCollector)
	if err != nil {
		logger.Errorw("Error creating object store",
			"error", err,
		)
	}

	kafkaConsumer, err := kafka.CreateKafkaConsumer(logger, config.Kafka.Brokers, config.Kafka.Topic)
	if err != nil {
		logger.Errorw("Error creating kafka consumer",
			"error", err,
		)
	}
	eventsManager, err := events.CreateEventsManager(logger, scanChan, kafkaConsumer, eventsCollector)
	if err != nil {
		logger.Errorw("Error creating events manager",
			"error", err,
		)
	}

	clamAV, err := clamav.CreateClamAV(logger)
	objectScanner, err := scanner.CreateObjectScanner(logger, objectStore, []scanner.Antivirus{clamAV}, auditLogger, scanCollector, config.ClamAV.RemoveAfterScan, config.ClamAV.DatetimeFormat, config.ClamAV.Path)
	dispatcher, err := dispatcher.CreateDispatcher(logger, []dispatcher.Scanner{objectScanner}, scanChan)
	if err != nil {
		logger.Errorw("Error creating dispatcher",
			"error", err,
		)
	}

	// sync.WaitGroup() as part of termination
	go eventsManager.Start()
	go dispatcher.Start()
	go metrics.Start()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan // Wait until interrupt
	logger.Infoln("Shutting down Aegis")
	// Cleanup stuff ...
	eventsManager.Stop()
	dispatcher.Stop()
	metrics.Stop()
	// Only stop when all scans finished?
	// Send signals to kafka and scan maanger to stop
	// sync.waitgroup to wait for all scans to finish
	return 0
}

func main() {
	os.Exit(run())
}
