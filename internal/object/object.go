package object

import (
	"aegis/internal/config"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Object struct {
	sugar      *zap.SugaredLogger
	ObjectKey  string
	BucketName string
	CachePath  string
	cachePerms fs.FileMode
	// Byte stream? Will that avoid saving to file?
}

func CreateObject(sugar *zap.SugaredLogger, bucketName string, objectKey string) (*Object, error) {
	// Make sure cache dir exists??
	config, err := config.GetConfig()
	if err != nil {
		sugar.Errorw("Failed to get config in object: ",
			"error", err,
		)
		return nil, err
	}
	cachePath := config.CachePath + bucketName + "/" + objectKey
	sugar.Debugw("Cache path: ",
		"cachePath", cachePath,
	)
	cachePerms := fs.FileMode(config.CachePerms)
	return &Object{
		sugar:      sugar,
		ObjectKey:  objectKey,
		BucketName: bucketName,
		CachePath:  cachePath,
		cachePerms: cachePerms,
	}, nil
}

func (o *Object) SaveByteStreamToFile(objectStream []byte) error {
	// Check if the parent directory of the file exists, and create it if it doesn't exist
	destDir := filepath.Dir(o.CachePath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, o.cachePerms); err != nil {
			o.sugar.Errorw("Failed to create parent directory: ",
				"error", err,
			)
			return err
		}
	}

	// Write the byte stream to the file
	err := ioutil.WriteFile(o.CachePath, objectStream, o.cachePerms)
	if err != nil {
		o.sugar.Errorw("Failed to save byte stream to file: ",
			"error", err,
		)
		return err
	}

	return nil
}

func (o *Object) RemoveFileFromCache() error {
	o.sugar.Errorw("Removing file from cache: ", o.CachePath)
	err := os.Remove(o.CachePath)
	if err != nil {
		o.sugar.Errorw("Failed to remove file from cache: ",
			"error", err,
		)
	}
	return nil
}
