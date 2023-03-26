package main

import (
	"antivirus/internal/kafka"
	"antivirus/internal/scan"
	"fmt"
	"os"
	"time"
)

func run() int {
	topic := "test-topic"
	scanChan := make(chan string)
	kafkaMgr, err := kafkaMgr.CreateKafkaMgr(topic, scanChan)
	if err != nil {
		fmt.Println("Error creating kafka manager")
	}

	scanMgr, err := scanMgr.CreateScanMgr(scanChan)
	if err != nil {
		fmt.Println("Error creating antivirus manager")
	}

	go kafkaMgr.StartKafkaMgr()
	go scanMgr.StartScanMgr()
	// go metricMgr.StartMetrics()

	// audit := StartAudit()
	// metric := StartMetrics()

	time.Sleep(10 * time.Minute)
	return 0
}

func shutdown() {
	fmt.Println("Shutting down")
}

func StartAudit() {
	fmt.Println("Starting Audit")
}

func StartMetrics() {
	fmt.Println("Starting Metrics")
}

func main() {
	os.Exit(run())
}
