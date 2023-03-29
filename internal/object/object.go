package object

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Object struct {
	ObjectKey  string
	BucketName string
	CachePath  string
	// Byte stream? Will that avoid saving to file?
}

func CreateObject(bucketName string, objectKey string) *Object {
	cachePath := "cache/" + bucketName + objectKey
	return &Object{ObjectKey: objectKey, BucketName: bucketName, CachePath: cachePath}
}

// Helper functions to do with objects e.g. save to file
func (o *Object) SaveByteStreamToFile(objectStream []byte) error {
	// TODO Check perms here
	err := ioutil.WriteFile(o.CachePath, objectStream, 0644)
	if err != nil {
		fmt.Println("Failed to save byte stream to file: ", err)
		return err
	}
	return nil
}

func (o *Object) RemoveFileFromCache() error {
	fmt.Println("Removing file from cache: ", o.CachePath)
	err := os.Remove(o.CachePath)
	if err != nil {
		fmt.Println("Failed to remove file from cache: ", err)
	}
	return nil
}

func (o *Object) GenerateHash() (string, error) {
	hash := "hash"
	return hash, nil
}
