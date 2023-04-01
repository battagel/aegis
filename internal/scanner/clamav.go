package scanner

import (
	"antivirus/internal/config"
	"antivirus/internal/object"
	"fmt"
	"os/exec"
	"time"
)

type ObjectStore interface {
	GetObject(bucketName string, objectName string) ([]byte, error)
	AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error
}

type Scanner struct {
	objectStore     ObjectStore
	removeAfterScan bool
	datetimeFormat  string
}

func CreateClamAV(objectStore ObjectStore) (*Scanner, error) {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config in clamav: ", err)
		return nil, err
	}
	removeAfterScan := config.Services.ClamAV.RemoveAfterScan
	datetimeFormat := config.Services.ClamAV.DatetimeFormat
	return &Scanner{objectStore: objectStore, removeAfterScan: removeAfterScan, datetimeFormat: datetimeFormat}, nil
}

func (s *Scanner) ScanObject(object *object.Object) error {

	objectStream, err := s.objectStore.GetObject(object.BucketName, object.ObjectKey)
	if err != nil {
		fmt.Println("Error getting object from object store: ", err)
		return err
	}

	err = object.SaveByteStreamToFile(objectStream)
	if err != nil {
		fmt.Println("Error saving byte stream to file: ", err)
		return err
	}

	result, err := s.executeScan(object.CachePath)
	if err != nil {
		fmt.Println("Error initiating scan: ", err)
		return err
	}
	dt := time.Now().Format(s.datetimeFormat)
	if result {
		fmt.Println("Infected file: ", object.ObjectKey)
		newTags := map[string]string{"antivirus": "infected", "antivirus-last-scanned": dt}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			fmt.Println("Error adding tag to object: ", err)
			return err
		}
	} else {
		fmt.Println("Clean file: ", object.ObjectKey)
		newTags := map[string]string{"antivirus": "scanned", "antivirus-last-scanned": dt}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			fmt.Println("Error adding tag to object: ", err)
			return err
		}
	}
	if s.removeAfterScan {
		err := object.RemoveFileFromCache()
		if err != nil {
			fmt.Println("Error removing file from cache: ", err)
			return err
		}
	}
	return nil
}

func (c *Scanner) executeScan(filePath string) (bool, error) {
	// Returns false if file is clean, true if infected
	// If there are any errors then return true (infected)
	cmd := exec.Command("clamdscan", filePath, "--stream", "-m")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means infected
			if exitError.ExitCode() == 1 {
				fmt.Println("File is infected")
				return true, nil
			}
		}
		fmt.Println("Error initiating scan: ", err)
		return true, err
	}
	// Due to exit codes, the file must be ok
	return false, nil
}
