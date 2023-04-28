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
	"aegis/pkg/postgresql"
	"aegis/pkg/prometheus"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config", err)
		return 1
	}
	cli.PrintSplash(config.LoggerEncoding)
	logger, err := logger.CreateZapLogger(config.LoggerLevel, config.LoggerEncoding)
	if err != nil {
		fmt.Println("Error creating logger", err)
		return 1
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
		return 1
	}
	objectStoreCollector, err := objectstore.CreateObjectStoreCollector(logger)
	if err != nil {
		logger.Errorw("Error creating object store collector",
			"error", err,
		)
		return 1
	}
	eventsCollector, err := events.CreateKafkaCollector(logger)
	if err != nil {
		logger.Errorw("Error creating kafka collector",
			"error", err,
		)
		return 1
	}
	scanCollector, err := scanner.CreateScanCollector(logger)
	if err != nil {
		logger.Errorw("Error creating collectors",
			"error", err,
		)
		return 1
	}

	postgresqlDB, dbClose, err := postgresql.CreatePostgresqlDB(logger, config.PostgresqlUsername, config.PostgresqlPassword, config.PostgresqlEndpoint, config.PostgresqlDatabase)
	if err != nil {
		logger.Errorw("Error creating postgresql database",
			"error", err,
		)
		return 1
	}
	defer dbClose()
	auditLogger, err := auditlog.CreateAuditLogger(logger, postgresqlDB, config.PostgresqlTable)
	if err != nil {
		logger.Errorw("Error creating audit logger",
			"error", err,
		)
		return 1
	}

	minioStore, err := minio.CreateMinio(logger, config.MinioEndpoint, config.MinioAccessKey, config.MinioSecretKey, config.MinioUseSSL)
	if err != nil {
		logger.Errorw("Error creating minio client",
			"error", err,
		)
		return 1
	}
	objectStore, err := objectstore.CreateObjectStore(logger, minioStore, objectStoreCollector)
	if err != nil {
		logger.Errorw("Error creating object store",
			"error", err,
		)
		return 1
	}

	kafkaConsumer, err := kafka.CreateKafkaConsumer(logger, config.KafkaBrokers, config.KafkaTopic)
	if err != nil {
		logger.Errorw("Error creating kafka consumer",
			"error", err,
		)
		return 1
	}
	eventsManager, err := events.CreateEventsManager(logger, scanChan, kafkaConsumer, eventsCollector)
	if err != nil {
		logger.Errorw("Error creating events manager",
			"error", err,
		)
		return 1
	}

	clamAV, err := clamav.CreateClamAV(logger)
	objectScanner, err := scanner.CreateObjectScanner(logger, objectStore, []scanner.Antivirus{clamAV}, auditLogger, scanCollector, config.ClamAVRemoveAfterScan, config.ClamAVDateTimeFormat, config.ClamAVPath)
	dispatcher, err := dispatcher.CreateDispatcher(logger, []dispatcher.Scanner{objectScanner}, scanChan)
	if err != nil {
		logger.Errorw("Error creating dispatcher",
			"error", err,
		)
		return 1
	}

	errChan := make(chan error)
	go eventsManager.Start(errChan)
	go dispatcher.Start()
	go metrics.Start(errChan)

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	select {
	case err = <-errChan:
		logger.Errorw("Error in goroutines",
			"error", err,
		)
	case <-sigchan:
		logger.Infoln("Shutting down Aegis")
	}
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
