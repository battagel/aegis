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
}

type ObjectStore struct {
	logger               logger.Logger
	minio                Minio
	objectStoreCollector ObjectStoreCollector
}

func CreateObjectStore(logger logger.Logger, minio Minio, objectStoreCollector ObjectStoreCollector) (*ObjectStore, error) {
	logger.Debugln("Creating object store")
	return &ObjectStore{
		logger:               logger,
		minio:                minio,
		objectStoreCollector: objectStoreCollector,
	}, nil
}

func (o *ObjectStore) GetObject(bucketName, objectKey string) ([]byte, error) {
	o.logger.Debugw("Getting object",
		"bucketName", bucketName,
		"objectKey", objectKey,
	)
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
	o.logger.Debugw("Putting object",
		"bucketName", bucketName,
		"objectKey", objectKey,
	)
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
	o.logger.Debugw("Removing object",
		"bucketName", bucketName,
		"objectKey", objectKey,
	)
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
	o.logger.Debugw("Moving object",
		"bucketName", bucketName,
		"objectKey", objectKey,
		"newBucketName", newBucketName,
		"newObjectKey", newObjectKey,
	)
	object, err := o.minio.GetObject(bucketName, objectKey)
	if err != nil {
		o.logger.Errorw("Error getting object from minio",
			"bucketName", bucketName,
			"objectKey", objectKey,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.GetObject()
	err = o.minio.PutObject(newBucketName, newObjectKey, object)
	if err != nil {
		o.logger.Errorw("Error putting object to minio",
			"bucketName", newBucketName,
			"objectKey", newObjectKey,
			"error", err,
		)
		return err
	}
	o.objectStoreCollector.PutObject()
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
	return nil
}

func (o *ObjectStore) GetObjectTagging(bucketName, objectKey string) (map[string]string, error) {
	o.logger.Debugw("Getting object tagging",
		"bucketName", bucketName,
		"objectKey", objectKey,
	)
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
	o.logger.Debugw("Putting object tagging",
		"bucketName", bucketName,
		"objectKey", objectKey,
		"tags", tags,
	)
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

func (o *ObjectStore) AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error {
	o.logger.Debugw("Adding object tag",
		"bucketName", bucketName,
		"objectName", objectName,
		"newTags", newTags,
	)
	objectTags, err := o.minio.GetObjectTagging(bucketName, objectName)
	if err != nil {
		o.logger.Errorw("Error getting object tags",
			"error", err,
			"bucketName", bucketName,
			"objectName", objectName,
			"newTags", newTags,
		)
		return err
	}
	o.objectStoreCollector.GetObjectTagging()
	for key, value := range newTags {
		objectTags[key] = value
	}
	err = o.minio.PutObjectTagging(bucketName, objectName, objectTags)
	if err != nil {
		o.logger.Errorw("Error setting object tag",
			"error", err,
			"bucketName", bucketName,
			"objectName", objectName,
			"newTags", newTags,
		)
		return err
	}
	o.objectStoreCollector.PutObjectTagging()
	return nil
}
