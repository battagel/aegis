package main

import (
	"aegis/internal/cli"
	"aegis/internal/dispatcher"
	"aegis/internal/kafka"
	"aegis/internal/metrics"
	"aegis/internal/object"
	"aegis/internal/objectstore"
	"aegis/internal/scanner"
	"log"
	"os"
	"os/signal"
)

func run() int {
	cli.PrintSplash()
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	scanChan := make(chan *object.Object)
	defer close(scanChan)
	metricChan := make(chan string)
	defer close(metricChan)

	// Removes hidden control flow
	logger.Println("Starting Aegis")
	metricManager, err := metrics.CreateMetricManager(metricChan)
	if err != nil {
		logger.Println("Error creating metric server")
	}
	objectStore, err := objectstore.CreateObjectStore()
	if err != nil {
		logger.Println("Error creating minio manager")
	}
	kafkaManager, err := kafka.CreateKafkaManager(scanChan)
	if err != nil {
		logger.Println("Error creating kafka manager")
	}
	clamAV, err := scanner.CreateClamAV(objectStore, metricChan)
	dispatcher, err := dispatcher.CreateDispatcher([]dispatcher.Scanner{clamAV}, scanChan)
	if err != nil {
		logger.Println("Error creating antivirus manager")
	}

	// sync.WaitGroup() as part of termination
	go kafkaManager.StartKafkaManager()
	go dispatcher.StartDispatcher()
	go metricManager.StartMetricManager()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan // Wait until interrupt
	logger.Println("Shutting down Aegis")
	// Cleanup stuff ...
	// Only stop when all scans finished?
	// Send signals to kafka and scan maanger to stop
	// sync.waitgroup to wait for all scans to finish
	return 0
}

func main() {
	os.Exit(run())
}
