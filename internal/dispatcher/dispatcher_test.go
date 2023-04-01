package dispatcher

import (
	"antivirus/internal/objectstore"
	"antivirus/internal/scanner"
	"fmt"
	"testing"
)

type CommonTestItems struct {
	scanners []Scanner
}

func ProvideCommonTestItems(t *testing.T) *CommonTestItems {
	t.Helper()

	objectStore, err := objectstore.CreateMinioManager()
	if err != nil {
		fmt.Println("Error creating minio manager")
	}
	clamAV, err := scanner.CreateClamAV(objectStore)
	if err != nil {
		fmt.Println("Error creating clamAV")
	}

	return &CommonTestItems{scanners: []Scanner{clamAV}}
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
