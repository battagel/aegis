package object

import (
	"aegis/internal/config"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

type Object struct {
	ObjectKey  string
	BucketName string
	CachePath  string
	cachePerms fs.FileMode
	// Byte stream? Will that avoid saving to file?
}

func CreateObject(bucketName string, objectKey string) *Object {
	// Make sure cache dir exists??
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Failed to get config in object: ", err)
	}
	cachePath := config.CachePath + bucketName + objectKey
	cachePerms := fs.FileMode(config.CachePerms)
	return &Object{ObjectKey: objectKey, BucketName: bucketName, CachePath: cachePath, cachePerms: cachePerms}
}

// Helper functions to do with objects e.g. save to file
func (o *Object) SaveByteStreamToFile(objectStream []byte) error {
	// TODO Check perms here
	err := ioutil.WriteFile(o.CachePath, objectStream, o.cachePerms)
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
	fmt.Println()
	return nil
}

func (o *Object) GenerateHash() (string, error) {
	hash := "hash"
	return hash, nil
}
