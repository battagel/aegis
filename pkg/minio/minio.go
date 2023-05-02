package minio

import (
	"aegis/pkg/logger"
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

type Minio struct {
	logger      logger.Logger
	minioClient *minio.Client
}

func CreateMinio(logger logger.Logger, endpoint, accessKey, secretKey string, useSSL bool) (*Minio, error) {
	logger.Debugln("Creating MinIO client")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Errorw("Connecting to MinIO failed",
			"error", err,
		)
	}
	return &Minio{
		logger:      logger,
		minioClient: minioClient,
	}, nil
}

func (m *Minio) GetObject(bucketName string, objectName string) ([]byte, error) {
	object, err := m.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		m.logger.Errorw("Error getting object",
			"error", err,
		)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		m.logger.Errorw("Error reading object",
			"error", err,
		)
		return nil, err
	}
	return data, nil
}

func (m *Minio) PutObject(bucketName string, objectName string, data []byte) error {
	_, err := m.minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		m.logger.Errorw("Error putting object",
			"error", err,
		)
		return err
	}
	return nil
}

func (m *Minio) RemoveObject(bucketName string, objectName string) error {
	err := m.minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		m.logger.Errorw("Error removing object",
			"error", err,
		)
		return err
	}
	return nil
}

func (m *Minio) GetObjectTagging(bucketName string, objectName string) (map[string]string, error) {
	tags, err := m.minioClient.GetObjectTagging(context.Background(), bucketName, objectName, minio.GetObjectTaggingOptions{})
	if err != nil {
		m.logger.Errorw("Error getting object tags",
			"error", err,
		)
		return nil, err
	}
	return tags.ToMap(), nil
}

func (m *Minio) PutObjectTagging(bucketName string, objectName string, newTags map[string]string) error {
	putTags, err := tags.MapToObjectTags(newTags)
	if err != nil {
		m.logger.Errorw("Error converting map to object tags",
			"error", err,
		)
		return err
	}
	err = m.minioClient.PutObjectTagging(context.Background(), bucketName, objectName, putTags, minio.PutObjectTaggingOptions{})
	if err != nil {
		m.logger.Errorw("Error setting object tag",
			"error", err,
		)
		return err
	}
	return nil
}

func (m *Minio) AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error {
	// Adds a tag by getting tags and ammending them
	objectTags, err := m.GetObjectTagging(bucketName, objectName)
	if err != nil {
		m.logger.Errorw("Error getting object tags",
			"error", err,
		)
		return err
	}
	for key, value := range newTags {
		objectTags[key] = value
	}
	err = m.PutObjectTagging(bucketName, objectName, objectTags)
	if err != nil {
		m.logger.Errorw("Error setting object tag",
			"error", err,
		)
		return err
	}
	return nil
}
