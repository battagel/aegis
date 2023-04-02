package main

import (
	"aegis/internal/dispatcher"
	"aegis/internal/kafka"
	"aegis/internal/metrics"
	"aegis/internal/object"
	"aegis/internal/objectstore"
	"aegis/internal/scanner"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	scanChan := make(chan *object.Object)
	metricChan := make(chan string)
	defer close(scanChan)
	// Removes hidden control flow
	objectStore, err := objectstore.CreateObjectStore()
	if err != nil {
		fmt.Println("Error creating minio manager")
	}
	kafkaManager, err := kafka.CreateKafkaManager(scanChan)
	if err != nil {
		fmt.Println("Error creating kafka manager")
	}
	clamAV, err := scanner.CreateClamAV(objectStore)
	dispatcher, err := dispatcher.CreateDispatcher([]dispatcher.Scanner{clamAV}, scanChan)
	if err != nil {
		fmt.Println("Error creating antivirus manager")
	}
	metricManager, err := metrics.CreateMetricManager(metricChan)
	if err != nil {
		fmt.Println("Error creating metric server")
	}

	// sync.WaitGroup() as part of termination
	go kafkaManager.StartKafkaManager()
	go dispatcher.StartDispatcher()
	go metricManager.StartMetricManager()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan // Wait until interrupt
	fmt.Println("Shutting down")
	// Cleanup stuff ...
	// Only stop when all scans finished?
	// Send signals to kafka and scan maanger to stop
	// sync.waitgroup to wait for all scans to finish
	return 0
}

func shutdown() {
	fmt.Println("Shutting down")
}

func main() {
	os.Exit(run())
}
