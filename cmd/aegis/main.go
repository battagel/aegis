package main

import (
	"aegis/internal/cli"
	"aegis/internal/dispatcher"
	"aegis/internal/kafka"
	"aegis/internal/metrics"
	"aegis/internal/object"
	"aegis/internal/objectstore"
	"aegis/internal/scanner"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func run() int {
	cli.PrintSplash()
	logger, err := zap.NewDevelopment()
	// logger, err := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	scanChan := make(chan *object.Object)
	defer close(scanChan)

	// Removes hidden control flow
	sugar.Infoln("Starting Aegis")
	metricManager, err := metrics.CreateMetricManager(sugar)
	if err != nil {
		sugar.Errorw("Error creating metric server",
			"error", err,
		)
	}
	objectStoreCollector, err := objectstore.CreateObjectStoreCollector(sugar)
	if err != nil {
		sugar.Errorw("Error creating object store collector",
			"error", err,
		)
	}
	objectStore, err := objectstore.CreateObjectStore(sugar, objectStoreCollector)
	if err != nil {
		sugar.Errorw("Error creating object store",
			"error", err,
		)
	}

	kafkaCollector, err := kafka.CreateKafkaCollector(sugar)
	if err != nil {
		sugar.Errorw("Error creating kafka collector",
			"error", err,
		)
	}
	kafkaManager, err := kafka.CreateKafkaManager(sugar, scanChan, kafkaCollector)
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
	// Only stop when all scans finished?
	// Send signals to kafka and scan maanger to stop
	// sync.waitgroup to wait for all scans to finish
	return 0
}

func main() {
	os.Exit(run())
}
