package main

import (
	"aegis/internal/auditlog"
	"aegis/internal/cleaner"
	"aegis/internal/cli"
	"aegis/internal/dispatcher"
	"aegis/internal/events"
	"aegis/internal/metrics"
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
	"context"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	// ### Initialisation and Configuration ###
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Infoln("Starting Aegis")
	logger.Infow("Config",
		"config", config,
	)
	logger.Debugln("Creating metric manager and collectors")
	prometheus, err := prometheus.CreatePrometheusExporter(logger, ctx, config.PrometheusEndpoint, config.PrometheusPath)
	if err != nil {
		logger.Errorw("Error creating prometheus server",
			"error", err,
		)
		return 1
	}
	metrics, err := metrics.CreateMetricsManager(logger, prometheus)
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
		logger.Errorw("Error creating scan collectors",
			"error", err,
		)
		return 1
	}
	cleanerCollector, err := cleaner.CreateCleanerCollector(logger)
	if err != nil {
		logger.Errorw("Error creating cleaner collectors",
			"error", err,
		)
		return 1
	}

	logger.Debugln("Creating audit logger")
	postgresqlDB, dbClose, err := postgresql.CreatePostgresqlDB(logger, ctx, config.PostgresqlUsername, config.PostgresqlPassword, config.PostgresqlEndpoint, config.PostgresqlDatabase)
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

	logger.Debugln("Creating object store")
	minioStore, err := minio.CreateMinio(logger, ctx, config.MinioEndpoint, config.MinioAccessKey, config.MinioSecretKey, config.MinioUseSSL)
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

	logger.Debugln("Creating events system")
	kafkaConsumer, err := kafka.CreateKafkaConsumer(logger, config.KafkaBrokers, config.KafkaTopic, config.KafkaGroupID, config.KafkaMaxBytes)
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

	logger.Debugln("Creating scanning workflow")
	clamAV, err := clamav.CreateClamAV(logger)
	if err != nil {
		logger.Errorw("Error creating clamav",
			"error", err,
		)
		return 1
	}
	cleaner, err := cleaner.CreateCleaner(logger, objectStore, config.CleanupPolicy, config.QuarantineBucket, cleanerCollector, auditLogger)
	if err != nil {
		logger.Errorw("Error creating cleaner",
			"error", err,
		)
		return 1
	}
	objectScanner, err := scanner.CreateObjectScanner(logger, objectStore, []scanner.Antivirus{clamAV}, cleaner, auditLogger, scanCollector, config.ClamAVRemoveAfterScan, config.ClamAVDateTimeFormat, config.ClamAVPath)
	if err != nil {
		logger.Errorw("Error creating object scanner",
			"error", err,
		)
		return 1
	}
	dispatcher, err := dispatcher.CreateDispatcher(logger, []dispatcher.Scanner{objectScanner}, scanChan)
	if err != nil {
		logger.Errorw("Error creating dispatcher",
			"error", err,
		)
		return 1
	}

	// ### Main Loop ###
	eventCtx, eventCancel := context.WithCancel(context.Background())
	defer eventCancel()
	errChan := make(chan error)
	doneChan := make(chan struct{})
	shutdownChan := make(chan os.Signal)
	go eventsManager.Start(eventCtx, errChan)
	go dispatcher.Start(errChan, doneChan)
	go metrics.Start()

	// ### Shutdown Sequence ###
	signal.Notify(shutdownChan, os.Interrupt)
	select {
	case err = <-errChan:
		logger.Errorw("Error in goroutines",
			"error", err,
		)
		eventCancel()
	case <-shutdownChan:
		logger.Infoln("Shutting down Aegis")
		eventCancel()
	}
	<-doneChan
	metrics.Stop()
	return 0
}

func main() {
	os.Exit(run())
}
