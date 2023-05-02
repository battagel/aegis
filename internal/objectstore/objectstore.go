package objectstore

import (
	"aegis/pkg/logger"
)

type ObjectStoreCollector interface {
	GetObject()
	PutObject()
	RemoveObject()
	GetObjectTagging()
	PutObjectTagging()
}

type Minio interface {
	GetObject(string, string) ([]byte, error)
	PutObject(string, string, []byte) error
	RemoveObject(string, string) error
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
			"objectKey", objectKey,
			"error", err,
		)
		return nil, err
	}
	o.objectStoreCollector.GetObject()
	return object, nil
}

func (o *ObjectStore) PutObject(bucketName, objectKey string, object []byte) error {
	err := o.minio.PutObject(bucketName, objectKey, object)
	if err != nil {
		o.logger.Errorw("Error putting object to minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.PutObject()
	return nil
}

func (o *ObjectStore) RemoveObject(bucketName, objectKey string) error {
	err := o.minio.RemoveObject(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error removing object from minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.RemoveObject()
	return nil
}

func (o *ObjectStore) MoveObject(bucketName, objectKey, newBucketName, newObjectKey string) error {
	object, err := o.minio.GetObject(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error getting object from minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
			"error", err,
		)
		return err
	}
	err = o.minio.PutObject(newBucketName, newObjectKey, object)
	if err != nil {
		o.logger.Errorw("Error putting object to minio",
			"bucketName", newBucketName,
			"objectKey", newObjectKey,
			"error", err,
		)
		return err
	}
	err = o.minio.RemoveObject(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error removing object from minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.RemoveObject()
	o.objectStoreCollector.PutObject()
	return nil
}

func (o *ObjectStore) GetObjectTagging(bucketName, objectKey string) (map[string]string, error) {
	tags, err := o.minio.GetObjectTagging(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error getting object tagging from minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
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
			"objectKey", objectKey,
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
			"objectKey", objectKey,
			"tags", tags,
			"error", err,
		)
		return err
	}
	return nil
}
