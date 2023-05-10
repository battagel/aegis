package scanner

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"time"
)

type ObjectStore interface {
	GetObject(string, string) ([]byte, error)
	AddObjectTagging(string, string, map[string]string) error
}

type Antivirus interface {
	ScanFile(string) (bool, string, error)
	GetName() string
}

type Cleaner interface {
	Cleanup(*object.Object, bool, string) error
}

type AuditLogger interface {
	Log(string, string, string, string, string, string)
}

type ScanCollector interface {
	FileScanned()
	CleanFile()
	InfectedFile()
	ScanError()
	ScanTime(float64)
}

type Scanner struct {
	logger          logger.Logger
	objectStore     ObjectStore
	antiviruses     []Antivirus
	cleaner         Cleaner
	auditLogger     AuditLogger
	scanCollector   ScanCollector
	removeAfterScan bool
	datetimeFormat  string
	cachePath       string
}

func CreateObjectScanner(logger logger.Logger, objectStore ObjectStore, antiviruses []Antivirus, cleaner Cleaner, auditLogger AuditLogger, scanCollector ScanCollector, removeAfterScan bool, datetimeFormat string, cachePath string) (*Scanner, error) {
	// Scanner for antiviruses that need the file downloaded
	return &Scanner{
		logger:          logger,
		objectStore:     objectStore,
		antiviruses:     antiviruses,
		cleaner:         cleaner,
		auditLogger:     auditLogger,
		scanCollector:   scanCollector,
		removeAfterScan: removeAfterScan,
		datetimeFormat:  datetimeFormat,
		cachePath:       cachePath,
	}, nil
}

func (s *Scanner) ScanObject(object *object.Object, errChan chan error) {
	scanTime := time.Now().Format(s.datetimeFormat)
	object.SetCachePath(s.cachePath)
	s.logger.Debugw("Getting object from object store",
		"bucketName", object.BucketName,
		"objectKey", object.ObjectKey,
	)
	objectStream, err := s.objectStore.GetObject(object.BucketName, object.ObjectKey)
	if err != nil {
		s.logger.Errorw("Error getting object from object store",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"error", err,
		)
		s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_getting_object", "", scanTime, "")
		s.scanCollector.ScanError()
		s.logger.Debugln("Debug")
		errChan <- err
		return
	}

	s.logger.Debugw("Saving byte stream to file",
		"bucketName", object.BucketName,
		"objectKey", object.ObjectKey,
	)
	err = object.SaveByteStreamToFile(objectStream)
	if err != nil {
		s.logger.Errorw("Error saving byte stream to file",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"error", err,
		)
		s.scanCollector.ScanError()
		errChan <- err
		return
	}

	s.logger.Debugw("Scanning file with antiviruses",
		"antiviruses", s.antiviruses,
		"bucketName", object.BucketName,
		"objectKey", object.ObjectKey,
	)
	overallResult := false
	for _, antivirus := range s.antiviruses {
		//scanStart := time.Now()
		result, virusType, err := antivirus.ScanFile(object.Path)
		//scanElapsed := float64(time.Since(scanStart) / time.Millisecond)
		// s.scanCollector.ScanTime(scanElapsed)
		if err != nil {
			s.logger.Errorw("Error executing scan",
				"antivirus", antivirus.GetName(),
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_scanning_file", antivirus.GetName(), scanTime, "")
			s.scanCollector.ScanError()
			errChan <- err
			return
		}
		s.scanCollector.FileScanned()
		// Seperate logging for each antivirus
		if result {
			overallResult = true
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "infected", antivirus.GetName(), scanTime, virusType)
		} else {
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "clean", antivirus.GetName(), scanTime, "")
		}
	}
	// Now process overall result
	if overallResult {
		s.logger.Warnw("Infected file",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		s.scanCollector.InfectedFile()
	} else {
		s.logger.Infow("Clean file",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		s.scanCollector.CleanFile()
	}
	s.cleaner.Cleanup(object, overallResult, scanTime)
	if s.removeAfterScan {
		err := object.RemoveFileFromCache()
		if err != nil {
			s.logger.Errorw("Error removing file from cache",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)

			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_removing_file_from_cache", "", scanTime, "")
			s.scanCollector.ScanError()
			errChan <- err
			return
		}
	}
	return
}
