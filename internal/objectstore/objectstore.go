package objectstore

import (
	"aegis/pkg/logger"
)

type ObjectStoreCollector interface {
	GetObject()
	GetObjectTagging()
	PutObjectTagging()
}

type Minio interface {
	GetObject(string, string) ([]byte, error)
	GetObjectTagging(string, string) (map[string]string, error)
	PutObjectTagging(string, string, map[string]string) error
	AddObjectTagging(string, string, map[string]string) error
}

type ObjectStore struct {
	logger               logger.Logger
	minio                Minio
	objectStoreCollector ObjectStoreCollector
}

func CreateObjectStore(logger logger.Logger, minio Minio, objectStoreCollector ObjectStoreCollector) (*ObjectStore, error) {
	logger.Debugln("Creating object store")
	return &ObjectStore{logger: logger, minio: minio, objectStoreCollector: objectStoreCollector}, nil
}

func (o *ObjectStore) GetObject(bucketName, objectKey string) ([]byte, error) {
	object, err := o.minio.GetObject(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error getting object from minio",
			"bucketName", bucketName,
			"objectName", objectKey,
			"error", err,
		)
		return nil, err
	}
	o.objectStoreCollector.GetObject()
	return object, nil
}

func (o *ObjectStore) GetObjectTagging(bucketName, objectKey string) (map[string]string, error) {
	tags, err := o.minio.GetObjectTagging(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error getting object tagging from minio",
			"bucketName", bucketName,
			"objectName", objectKey,
			"error", err,
		)
		return nil, err
	}
	o.objectStoreCollector.GetObjectTagging()
	return tags, nil
}

func (o *ObjectStore) PutObjectTagging(bucketName, objectKey string, tags map[string]string) error {
	err := o.minio.PutObjectTagging(bucketName, objectKey, tags)
	if err != nil {
		o.logger.Errorw("Error putting object tagging to minio",
			"bucketName", bucketName,
			"objectName", objectKey,
			"tags", tags,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.PutObjectTagging()
	return nil
}

func (o *ObjectStore) AddObjectTagging(bucketName, objectKey string, tags map[string]string) error {
	err := o.minio.AddObjectTagging(bucketName, objectKey, tags)
	if err != nil {
		o.logger.Errorw("Error adding object tagging to minio",
			"bucketName", bucketName,
			"objectName", objectKey,
			"tags", tags,
			"error", err,
		)
		return err
	}
	return nil
}
