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
	if err != nil {
		fmt.Println("Error getting config", err)
		return 0
	}
	logger, err := logger.CreateZapLogger(config.LoggerLevel, config.LoggerEncoding)
	if err != nil {
		fmt.Println("Error creating logger", err)
		return 0
	}

	scanChan := make(chan *object.Object)
	defer close(scanChan)

	logger.Infoln("Starting Aegis")
	logger.Infow("Config",
		"config", config,
	)
	logger.Debugln("Creating metric collectors")
	metrics, err := prometheus.CreatePrometheusServer(logger, config.PrometheusEndpoint, config.PrometheusPath)
	if err != nil {
		logger.Errorw("Error creating metric collectors",
			"error", err,
		)
		return 0
	}
	objectStoreCollector, err := objectstore.CreateObjectStoreCollector(logger)
	if err != nil {
		logger.Errorw("Error creating object store collector",
			"error", err,
		)
		return 0
	}
	eventsCollector, err := events.CreateKafkaCollector(logger)
	if err != nil {
		logger.Errorw("Error creating kafka collector",
			"error", err,
		)
		return 0
	}
	scanCollector, err := scanner.CreateScanCollector(logger)
	if err != nil {
		logger.Errorw("Error creating collectors",
			"error", err,
		)
		return 0
	}

	postgresDB, dbClose, err := postgres.CreatePostgresDB(logger, config.PostgresUsername, config.PostgresPassword, config.PostgresEndpoint, config.PostgresDatabase)
	if err != nil {
		logger.Errorw("Error creating postgres database",
			"error", err,
		)
		return 0
	}
	defer dbClose()
	auditLogger, err := auditlog.CreateAuditLogger(logger, postgresDB, config.PostgresTable)

	minioStore, err := minio.CreateMinio(logger, config.MinioEndpoint, config.MinioAccessKey, config.MinioSecretKey, config.MinioUseSSL)
	if err != nil {
		logger.Errorw("Error creating minio client",
			"error", err,
		)
		return 0
	}
	objectStore, err := objectstore.CreateObjectStore(logger, minioStore, objectStoreCollector)
	if err != nil {
		logger.Errorw("Error creating object store",
			"error", err,
		)
		return 0
	}

	kafkaConsumer, err := kafka.CreateKafkaConsumer(logger, config.KafkaBrokers, config.KafkaTopic)
	if err != nil {
		logger.Errorw("Error creating kafka consumer",
			"error", err,
		)
		return 0
	}
	eventsManager, err := events.CreateEventsManager(logger, scanChan, kafkaConsumer, eventsCollector)
	if err != nil {
		logger.Errorw("Error creating events manager",
			"error", err,
		)
		return 0
	}

	clamAV, err := clamav.CreateClamAV(logger)
	objectScanner, err := scanner.CreateObjectScanner(logger, objectStore, []scanner.Antivirus{clamAV}, auditLogger, scanCollector, config.ClamAVRemoveAfterScan, config.ClamAVDateTimeFormat, config.ClamAVPath)
	dispatcher, err := dispatcher.CreateDispatcher(logger, []dispatcher.Scanner{objectScanner}, scanChan)
	if err != nil {
		logger.Errorw("Error creating dispatcher",
			"error", err,
		)
		return 0
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
