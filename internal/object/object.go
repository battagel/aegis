package object

import (
	"aegis/internal/config"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
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
	cachePath := config.CachePath + bucketName + "/" + objectKey
	cachePerms := fs.FileMode(config.CachePerms)
	return &Object{ObjectKey: objectKey, BucketName: bucketName, CachePath: cachePath, cachePerms: cachePerms}
}

func (o *Object) SaveByteStreamToFile(objectStream []byte) error {
	// Check if the parent directory of the file exists, and create it if it doesn't exist
	destDir := filepath.Dir(o.CachePath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, o.cachePerms); err != nil {
			fmt.Println("Failed to create parent directory: ", err)
			return err
		}
	}

	// Write the byte stream to the file
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
