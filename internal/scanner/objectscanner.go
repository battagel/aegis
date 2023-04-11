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
	ScanFile(filePath string) (bool, error)
	GetName() string
}

type AuditLogger interface {
	Log(string, string, string, string, string)
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
	auditLogger     AuditLogger
	scanCollector   ScanCollector
	removeAfterScan bool
	datetimeFormat  string
	cachePath       string
}

func CreateObjectScanner(logger logger.Logger, objectStore ObjectStore, antiviruses []Antivirus, auditLogger AuditLogger, scanCollector ScanCollector, removeAfterScan bool, datetimeFormat string, cachePath string) (*Scanner, error) {
	// Scanner for antiviruses that need the file downloaded
	return &Scanner{
		logger:          logger,
		objectStore:     objectStore,
		antiviruses:     antiviruses,
		auditLogger:     auditLogger,
		scanCollector:   scanCollector,
		removeAfterScan: removeAfterScan,
		datetimeFormat:  datetimeFormat,
		cachePath:       cachePath,
	}, nil
}

func (s *Scanner) ScanObject(object *object.Object) error {
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
		s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_getting_object", "", scanTime)
		s.scanCollector.ScanError()
		return err
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
		return err
	}

	s.logger.Debugw("Scanning file with antiviruses",
		"antiviruses", s.antiviruses,
		"bucketName", object.BucketName,
		"objectKey", object.ObjectKey,
	)
	overallResult := false
	for _, antivirus := range s.antiviruses {
		scanStart := time.Now()
		result, err := antivirus.ScanFile(object.Path)
		scanElapsed := float64(time.Since(scanStart) / time.Millisecond)
		s.scanCollector.ScanTime(scanElapsed)
		if err != nil {
			s.logger.Errorw("Error executing scan",
				"antivirus", antivirus.GetName(),
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_scanning_file", antivirus.GetName(), scanTime)
			s.scanCollector.ScanError()
			return err
		}
		s.scanCollector.FileScanned()
		if result {
			overallResult = true
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "infected", antivirus.GetName(), scanTime)
		} else {
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "clean", antivirus.GetName(), scanTime)
		}
	}
	if overallResult {
		s.scanCollector.InfectedFile()
		s.logger.Warnw("Infected file",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		newTags := map[string]string{"antivirus": "infected", "antivirus-last-scanned": scanTime}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			s.logger.Errorw("Error adding tag to object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_adding_tags", "", scanTime)
			s.scanCollector.ScanError()
			return err
		}
	} else {
		s.logger.Infow("Clean file",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		s.scanCollector.CleanFile()
		newTags := map[string]string{"antivirus": "scanned", "antivirus-last-scanned": scanTime}
		err := s.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			s.logger.Errorw("Error adding tag to object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_adding_tags", "", scanTime)
			s.scanCollector.ScanError()
			return err
		}
	}
	if s.removeAfterScan {
		err := object.RemoveFileFromCache()
		if err != nil {
			s.logger.Errorw("Error removing file from cache",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			s.auditLogger.Log(object.BucketName, object.ObjectKey, "error_removing_file_from_cache", "", scanTime)
			s.scanCollector.ScanError()
			return err
		}
	}
	return nil
}
