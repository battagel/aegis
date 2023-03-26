package scanMgr

import (
	minioMgr "antivirus/internal/minio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type ScanMgr struct {
	removeAfterScan bool
	minioMgr        *minioMgr.MinioMgr
	scanChan        chan string
}

func CreateScanMgr(scanChan chan string) (*ScanMgr, error) {
	fmt.Println("Creating Scan Manager")
	removeAfterScan := true

	minioMgr, err := minioMgr.CreateMinioMgr()
	if err != nil {
		return nil, err
	}
	return &ScanMgr{removeAfterScan: removeAfterScan, scanChan: scanChan, minioMgr: minioMgr}, nil
}

func (s *ScanMgr) StartScanMgr() error {
	fmt.Println("Starting Scan Manager")
	for {
		objectPath := <-s.scanChan

		v := strings.Split(objectPath, "/")
		bucketName := v[0]
		objectName := v[1]

		object, err := s.minioMgr.GetObject(bucketName, objectName)

		cachePath := "cache/" + objectName
		err = saveByteStreamToFile(object, cachePath)
		if err != nil {
			fmt.Println("Error saving byte stream to file: ", err)
			return err
		}

		result, err := s.initiateScan(cachePath)
		if err != nil {
			fmt.Println("Error initiating scan: ", err)
			return err
		}
		// TODO Initiate scan of file
		fmt.Println("Scan result: ", result)

		if s.removeAfterScan {
			err := removeFileFromCache(cachePath)
			if err != nil {
				fmt.Println("Error removing file from cache: ", err)
				return err
			}
		}
	}
}

func (s *ScanMgr) StopScanMgr() error {
	fmt.Println("Stopping Scan Manager")

	return nil
}

func (s *ScanMgr) initiateScan(filePath string) (bool, error) {
	// TODO Connect to clamd and initiate scan
	cmd := exec.Command("clamdscan", filePath, "--stream", "-m")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error initiating scan: ", err)
		return true, err
	}
	fmt.Println("Scan output: ", string(output))

	return true, nil
}

func saveByteStreamToFile(data []byte, filePath string) error {
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Failed to save byte stream to file: ", err)
		return err
	}
	return nil
}

func removeFileFromCache(filePath string) error {
	fmt.Println("Removing file from cache: ", filePath)
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Failed to remove file from cache: ", err)
	}
	return nil
}
