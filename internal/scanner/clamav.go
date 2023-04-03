package scanner

import (
	"aegis/internal/config"
	"aegis/internal/object"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

type ObjectStore interface {
	GetObject(bucketName string, objectName string) ([]byte, error)
	AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error
}

type ScanCollector interface {
	FileScanned()
	CleanFile()
	InfectedFile()
	ScanError()
}

type Scanner struct {
	sugar           *zap.SugaredLogger
	objectStore     ObjectStore
	scanCollector   ScanCollector
	removeAfterScan bool
	datetimeFormat  string
}

func CreateClamAV(sugar *zap.SugaredLogger, objectStore ObjectStore, scanCollector ScanCollector) (*Scanner, error) {
	config, err := config.GetConfig()
	if err != nil {
		sugar.Errorw("Error getting config in clamav: ",
			"error", err,
		)
		return nil, err
	}
	removeAfterScan := config.Services.ClamAV.RemoveAfterScan
	datetimeFormat := config.Services.ClamAV.DatetimeFormat
	return &Scanner{
		sugar:           sugar,
		objectStore:     objectStore,
		scanCollector:   scanCollector,
		removeAfterScan: removeAfterScan,
		datetimeFormat:  datetimeFormat,
	}, nil
}

func (s *Scanner) ScanObject(object *object.Object) error {

	objectStream, err := s.objectStore.GetObject(object.BucketName, object.ObjectKey)
	if err != nil {
		s.sugar.Errorw("Error getting object from object store: ",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"error", err,
		)
		s.scanCollector.ScanError()
		return err
	}

	err = object.SaveByteStreamToFile(objectStream)
	if err != nil {
		s.sugar.Errorw("Error saving byte stream to file: ",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"error", err,
		)
		s.scanCollector.ScanError()
		return err
	}

	result, err := s.executeScan(object.CachePath)
	if err != nil {
		s.sugar.Errorw("Error executing scan: ",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"error", err,
		)
		s.scanCollector.ScanError()
		return err
	}
	s.scanCollector.FileScanned()
	dt := time.Now().Format(s.datetimeFormat)
	if result {
		s.sugar.Warnw("Infected file: ",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		s.scanCollector.InfectedFile()
		newTags := map[string]string{"antivirus": "infected", "antivirus-last-scanned": dt}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			s.sugar.Errorw("Error adding tag to object: ",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.scanCollector.ScanError()
			return err
		}
	} else {
		s.sugar.Infow("Clean file: ",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		s.scanCollector.CleanFile()
		newTags := map[string]string{"antivirus": "scanned", "antivirus-last-scanned": dt}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			s.sugar.Errorw("Error adding tag to object: ",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.scanCollector.ScanError()
			return err
		}
	}
	if s.removeAfterScan {
		err := object.RemoveFileFromCache()
		if err != nil {
			s.sugar.Errorw("Error removing file from cache: ",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.scanCollector.ScanError()
			return err
		}
	}
	return nil
}

func (s *Scanner) executeScan(filePath string) (bool, error) {
	// Returns false if file is clean, true if infected
	// If there are any errors then return true (infected)
	cmd := exec.Command("clamdscan", filePath, "--stream", "-m")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means infected
			if exitError.ExitCode() == 1 {
				s.sugar.Errorw("File is infected")
				return true, nil
			}
		}
		s.sugar.Errorw("Error running clamdscan. Is clamd running?",
			"error", err,
		)
		return true, err
	}
	// Due to exit codes, the file must be ok
	return false, nil
}
