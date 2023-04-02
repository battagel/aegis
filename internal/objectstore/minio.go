package objectstore

import (
	"aegis/internal/config"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

type ObjectStore struct {
	minioClient *minio.Client
}

func CreateObjectStore() (*ObjectStore, error) {
	fmt.Println("Creating object store")
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config in minio: ", err)
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
		fmt.Println("Connecting to MinIO failed: ", err)
		return nil, err
	}
	return &ObjectStore{minioClient: minioClient}, nil
}

func (m *ObjectStore) GetObject(bucketName string, objectName string) ([]byte, error) {
	object, err := m.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object: ", err)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		fmt.Println("Error reading object: ", err)
		return nil, err
	}
	return data, nil
}

func (m *ObjectStore) GetObjectTagging(bucketName string, objectName string) (*tags.Tags, error) {
	tags, err := m.minioClient.GetObjectTagging(context.Background(), bucketName, objectName, minio.GetObjectTaggingOptions{})
	if err != nil {
		fmt.Println("Error getting object tags: ", err)
		return nil, err
	}
	return tags, nil
}

func (m *ObjectStore) PutObjectTagging(bucketName string, objectName string, tags *tags.Tags) error {
	err := m.minioClient.PutObjectTagging(context.Background(), bucketName, objectName, tags, minio.PutObjectTaggingOptions{})
	if err != nil {
		fmt.Println("Error setting object tag: ", err)
		return err
	}
	return nil
}

func (m *ObjectStore) AddObjectTagging(bucketName string, objectName string, newTags map[string]string) error {
	// Adds a tag by getting tags and ammending them
	objectTags, err := m.GetObjectTagging(bucketName, objectName)
	if err != nil {
		fmt.Println("Error getting object tags: ", err)
		return err
	}
	tagMap := objectTags.ToMap()
	for key, value := range newTags {
		tagMap[key] = value
	}
	newObjectTags, err := tags.MapToObjectTags(tagMap)
	err = m.PutObjectTagging(bucketName, objectName, newObjectTags)
	if err != nil {
		fmt.Println("Error setting object tag: ", err)
		return err
	}
	return nil
}
