package scanner

import (
	"antivirus/internal/object"
	"fmt"
	"os/exec"
	// "github.com/minio/minio-go/v7/pkg/tags"
)

type ObjectStore interface {
	GetObject(bucketName string, objectName string) ([]byte, error)
	// AddObjectTagging(bucketName string, objectName string, newTags *tags.Tags) error
}

type Scanner struct {
	objectStore     ObjectStore
	removeAfterScan bool
}

func CreateClamAV(objectStore ObjectStore) (*Scanner, error) {
	removeAfterScan := true
	return &Scanner{objectStore: objectStore, removeAfterScan: removeAfterScan}, nil
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
	if result {
		// TODO: Send result to minio and add tag to file
		fmt.Println("Infected file: ", object.ObjectKey)
		// antivirusTags := AntivirusTags{"antivirus": "infected"}
		// err := s.minioMgr.AddObjectTagging(bucketName, objectKey, antivirusTags)
		// if err != nil {
		// fmt.Println("Error adding tag to object: ", err)
		// return err
		// }
	} else {
		fmt.Println("Clean file: ", object.ObjectKey)
		// antivirusTags := AntivirusTags{"antivirus": "scanned"}
		// err := s.minioMgr.AddObjectTagging(bucketName, objectKey, antivirusTags)
		// if err != nil {
		// 	fmt.Println("Error adding tag to object: ", err)
		// 	return err
		// }
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
