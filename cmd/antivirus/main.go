package main

import (
	"antivirus/internal/dispatcher"
	"antivirus/internal/kafka"
	"antivirus/internal/object"
	"antivirus/internal/objectstore"
	"antivirus/internal/scanner"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	topic := "test-topic"
	scanChan := make(chan *object.Object)
	// Removes hidden control flow
	objectStore, err := objectstore.CreateObjectStore()
	if err != nil {
		fmt.Println("Error creating minio manager")
	}
	kafkaManager, err := kafka.CreateKafkaManager(topic, scanChan)
	if err != nil {
		fmt.Println("Error creating kafka manager")
	}
	clamAV, err := scanner.CreateClamAV(objectStore)
	dispatcher, err := dispatcher.CreateDispatcher([]dispatcher.Scanner{clamAV}, scanChan)
	if err != nil {
		fmt.Println("Error creating antivirus manager")
	}

	// sync.WaitGroup() as part of termination
	go kafkaManager.StartKafkaManager()
	go dispatcher.StartDispatcher()

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
