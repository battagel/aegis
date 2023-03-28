package minio

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	// "github.com/minio/minio-go/v7/pkg/tags"
)

type MinioMgr struct {
	minioClient *minio.Client
}

func CreateMinioManager() (*MinioMgr, error) {
	// Must be an IP not http:// or https://
	endpoint := "192.168.1.138:9000"
	endpoint = "10.10.6.233:9000" // Work IP
	accessKey := "minioadmin"
	secretKey := "minioadmin"
	useSSL := false
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("Connecting to MinIO failed: ", err)
		return nil, err
	}
	return &MinioMgr{minioClient: minioClient}, nil
}

func (m *MinioMgr) GetObject(bucketName string, objectName string) ([]byte, error) {
	object, err := m.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object: ", err)
		return nil, err
	}
	defer object.Close()

	data, err := ioutil.ReadAll(object)
	if err != nil {
		fmt.Println("Error reading object: ", err)
	}
	return data, nil
}

// func (m *MinioMgr) GetObjectTagging(bucketName string, objectName string) (*tags.Tags, error) {
// tags, err := m.minioClient.GetObjectTagging(context.Background(), bucketName, objectName, ops)
// if err != nil {
// fmt.Println("Error getting object tags: ", err)
// return nil, err
// }
// return tags, nil
// }

// func (m *MinioMgr) SetObjectTagging(bucketName string, objectName string, tags *tags.Tags) error {
// err := m.minioClient.SetObjectTagging(context.Background(), bucketName, objectName, tags)
// if err != nil {
// fmt.Println("Error setting object tag: ", err)
// return err
// }
// return nil
// }
//
// func (m *MinioMgr) AddObjectTagging(bucketName string, objectName string, antivirusTags AntivirusTags) error {
// // Adds a tag by getting tags and ammending them
// tags, err := m.GetObjectTagging(bucketName, objectName)
// if err != nil {
// fmt.Println("Error getting object tags: ", err)
// return err
// }
// for tagName, tagValue := range antivirusTags {
// tags[tagName] = tagValue
// }
// ops := &minio.PutObjectTaggingOptions{}
// err := m.minioClient.PutObjectTagging(context.Background(), bucketName, objectName, tags)
// if err != nil {
// fmt.Println("Error setting object tag: ", err)
// return err
// }
// return nil
// }
