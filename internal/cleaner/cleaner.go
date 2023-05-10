package cleaner

import (
	"aegis/internal/object"
	"aegis/pkg/logger"
	"errors"
)

type ObjectStore interface {
	RemoveObject(string, string) error
	AddObjectTagging(string, string, map[string]string) error
	MoveObject(string, string, string, string) error
}

type AuditLogger interface {
	Log(string, string, string, string, string, string)
}

type CleanerCollector interface {
	ObjectTagged()
	ObjectRemoved()
	ObjectQuarantined()
	CleanupError()
}

type Cleaner struct {
	logger           logger.Logger
	objectStore      ObjectStore
	cleanupPolicy    string
	quarantineBucket string
	cleanerCollector CleanerCollector
	auditLogger      AuditLogger
}

func CreateCleaner(logger logger.Logger, objectStore ObjectStore, cleanupPolicy, quarantineBucket string, cleanerCollector CleanerCollector, auditLogger AuditLogger) (*Cleaner, error) {
	return &Cleaner{
		logger:           logger,
		objectStore:      objectStore,
		cleanupPolicy:    cleanupPolicy,
		quarantineBucket: quarantineBucket,
		cleanerCollector: cleanerCollector,
		auditLogger:      auditLogger,
	}, nil
}

func (c *Cleaner) Cleanup(object *object.Object, result bool, scanTime string) error {
	c.logger.Debugw("Cleaning up",
		"cleanupPolicy", c.cleanupPolicy,
	)
	var err error
	switch c.cleanupPolicy {
	case "tag":
		err = c.tagInfected(object, result, scanTime)
	case "remove":
		err = c.removeInfected(object, result, scanTime)
	case "quarantine":
		err = c.quarantineInfected(object, result, scanTime)
	default:
		c.logger.Warnln("No cleanup policy found")
	}
	if err != nil {
		c.logger.Errorw("Error cleaning up",
			"cleanupPolicy", c.cleanupPolicy,
			"error", err,
		)
		return err
	}
	return nil
}

func (c *Cleaner) tagInfected(object *object.Object, result bool, scanTime string) error {
	if result {
		c.logger.Debugw("Tagging infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		newTags := map[string]string{"antivirus": "infected", "antivirus-last-scanned": scanTime}
		err := c.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			c.logger.Errorw("Error adding tag to object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			c.auditLogger.Log(object.BucketName, object.ObjectKey, "error_adding_tags", "", scanTime, "")
			c.cleanerCollector.CleanupError()
			return err
		}
		c.logger.Debugw("Successfully tagged infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		c.cleanerCollector.ObjectTagged()
	} else {
		c.logger.Debugw("Tagging clean object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		newTags := map[string]string{"antivirus": "clean", "antivirus-last-scanned": scanTime}
		err := c.objectStore.AddObjectTagging(object.BucketName, object.ObjectKey, newTags)
		if err != nil {
			c.logger.Errorw("Error adding tag to object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			c.auditLogger.Log(object.BucketName, object.ObjectKey, "error_adding_tags", "", scanTime, "")
			c.cleanerCollector.CleanupError()
			return err
		}
		c.logger.Debugw("Successfully tagged clean object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		c.cleanerCollector.ObjectTagged()
	}
	return nil
}

func (c *Cleaner) removeInfected(object *object.Object, result bool, scanTime string) error {
	if result {
		c.logger.Debugw("Removing infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		err := c.objectStore.RemoveObject(object.BucketName, object.ObjectKey)
		if err != nil {
			c.logger.Errorw("Failed to remove infected object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"error", err,
			)
			c.auditLogger.Log(object.BucketName, object.ObjectKey, "error_removing_object", "", scanTime, "")
			c.cleanerCollector.CleanupError()
			return err
		}
		c.logger.Debugw("Successfully removed infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		c.cleanerCollector.ObjectRemoved()
	} else {
		c.logger.Debugw("Clean object, not removing",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
	}
	return nil
}

func (c *Cleaner) quarantineInfected(object *object.Object, result bool, scanTime string) error {
	if result {
		c.logger.Debugw("Quarantining infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
		if c.quarantineBucket == "" {
			c.logger.Errorw("Quarantine bucket not set",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
			)
			c.auditLogger.Log(object.BucketName, object.ObjectKey, "error_quarantining_object", "", scanTime, "")
			c.cleanerCollector.CleanupError()
			return errors.New("Quarantine bucket not set")
		}
		err := c.objectStore.MoveObject(object.BucketName, object.ObjectKey, c.quarantineBucket, object.ObjectKey)
		if err != nil {
			c.logger.Errorw("Failed to quarantine infected object",
				"bucketName", object.BucketName,
				"objectKey", object.ObjectKey,
				"quarantineBucket", c.quarantineBucket,
				"error", err,
			)
			c.auditLogger.Log(object.BucketName, object.ObjectKey, "error_quarantining_object", "", scanTime, "")
			c.cleanerCollector.CleanupError()
			return err
		}
		c.logger.Debugw("Successfully quarantined infected object",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
			"quarantineBucket", c.quarantineBucket,
		)
		c.cleanerCollector.ObjectQuarantined()
	} else {
		c.logger.Debugw("Clean object, not quarantining",
			"bucketName", object.BucketName,
			"objectKey", object.ObjectKey,
		)
	}
	return nil
}

// TODO Room for expansion... display infected files differently?
