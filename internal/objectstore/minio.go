package objectstore

import (
	"aegis/internal/config"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
	"go.uber.org/zap"
)

type ObjectStoreCollector interface {
	GetObject()
	GetObjectTagging()
	PutObjectTagging()
}

type ObjectStore struct {
	sugar                *zap.SugaredLogger
	minioClient          *minio.Client
	objectStoreCollector ObjectStoreCollector
}

func CreateObjectStore(sugar *zap.SugaredLogger, objectStoreCollector ObjectStoreCollector) (*ObjectStore, error) {
	sugar.Debugln("Creating object store")
	config, err := config.GetConfig()
	if err != nil {
		sugar.Errorw("Error getting config in minio",
			"error", err,
		)
		return nil, err
	}
	endpoint := config.Services.Minio.Endpoint
	accessKey := config.Services.Minio.AccessKey
	secretKey := config.Services.Minio.SecretKey
	useSSL := config.Services.Minio.UseSSL
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		sugar.Errorw("Connecting to MinIO failed",
			"error", err,
		)
		return nil, err
	}
	return &ObjectStore{sugar: sugar, minioClient: minioClient, objectStoreCollector: objectStoreCollector}, nil
}

func (m *ObjectStore) GetObject(bucketName string, objectName string) ([]byte, error) {
	object, err := m.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		m.sugar.Errorw("Error getting object",
			"error", err,
		)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		m.sugar.Errorw("Error reading object",
			"error", err,
		)
		return nil, err
	}
	m.objectStoreCollector.GetObject()
	return data, nil
}

func (m *ObjectStore) GetObjectTagging(bucketName string, objectName string) (*tags.Tags, error) {
	tags, err := m.minioClient.GetObjectTagging(context.Background(), bucketName, objectName, minio.GetObjectTaggingOptions{})
	if err != nil {
		m.sugar.Errorw("Error getting object tags",
			"error", err,
		)
		return nil, err
	}
	m.objectStoreCollector.GetObjectTagging()
	return tags, nil
}

func (m *ObjectStore) PutObjectTagging(bucketName string, objectName string, tags *tags.Tags) error {
	err := m.minioClient.PutObjectTagging(context.Background(), bucketName, objectName, tags, minio.PutObjectTaggingOptions{})
	if err != nil {
		m.sugar.Errorw("Error setting object tag",
			"error", err,
		)
		return err
	}
	m.objectStoreCollector.PutObjectTagging()
	return nil
}

func (m *ObjectStore) AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error {
	// Adds a tag by getting tags and ammending them
	objectTags, err := m.GetObjectTagging(bucketName, objectName)
	if err != nil {
		m.sugar.Errorw("Error getting object tags",
			"error", err,
		)
		return err
	}
	tagMap := objectTags.ToMap()
	for key, value := range newTags {
		tagMap[key] = value
	}
	newObjectTags, err := tags.MapToObjectTags(tagMap)
	err = m.PutObjectTagging(bucketName, objectName, newObjectTags)
	if err != nil {
		m.sugar.Errorw("Error setting object tag",
			"error", err,
		)
		return err
	}
	return nil
}
