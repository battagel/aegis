package dispatcher

import (
	"antivirus/internal/minioMgr"
	"fmt"
	"testing"
)

type CommonTestItems struct {
	ScanMgr *ScanMgr
}

func ProvideCommonTestItems(t *testing.T) *CommonTestItems {
	t.Helper()

	scanChan := make(chan string)
	minioMgr, err := minioMgr.CreateMinioMgr()
	if err != nil {
		fmt.Println("Error creating minio manager")
	}
	scanMgr, err := scanMgr.CreateScanMgr(scanChan, minioMgr)
	if err != nil {
		fmt.Println("Error creating antivirus manager")
	}

	go scanMgr.StartScanMgr()

	return &CommonTestItems{ScanMgr: scanMgr}
}

func TestScanMgrWorking(t *testing.T) {
	fmt.Println("TestScanMgrWorking")
	// Test the scan result of a non-existent file
	//
	// Put a dummy request into kafka queue

	// Test the scan result of the dummy request
}

func TestScanMgrFileNotFound(t *testing.T) {
	// Test the scan result of a non-existent file
}
