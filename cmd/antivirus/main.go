package main

import (
	"antivirus/internal/kafka"
	"antivirus/internal/messages"
	"antivirus/internal/minio"
	"antivirus/internal/scanmanager"
	"fmt"
	"os"
	"time"
)

func run() int {
	topic := "test-topic"
	scanChan := make(chan messages.ScanRequest)
	minioManager, err := minio.CreateMinioManager()
	if err != nil {
		fmt.Println("Error creating minio manager")
	}
	kafkaManager, err := kafka.CreateKafkaManager(topic, scanChan)
	if err != nil {
		fmt.Println("Error creating kafka manager")
	}

	scanManager, err := scanmanager.CreateScanManager(scanChan, minioManager)
	if err != nil {
		fmt.Println("Error creating antivirus manager")
	}

	go kafkaManager.StartKafkaManager()
	go scanManager.StartScanManager()
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
