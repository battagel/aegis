package minioMgr

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioMgr struct {
	minioClient *minio.Client
}

func CreateMinioMgr() (*MinioMgr, error) {
	endpoint := "192.168.1.138:9000"
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
	fmt.Println("Getting Minio Object")
	object, err := m.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object: ", err)
		return nil, err
	}
	fmt.Println(object)
	defer object.Close()

	data, err := ioutil.ReadAll(object)
	if err != nil {
		fmt.Println("Error reading object: ", err)
	}
	return data, nil
}
