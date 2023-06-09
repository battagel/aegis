package object

import (
	"aegis/pkg/logger"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Object struct {
	logger     logger.Logger
	ObjectKey  string
	BucketName string
	Perms      fs.FileMode
	Path       string
	// TODO Byte stream? Will that avoid saving to file?
}

const (
	// Permissions of the cache and subsequent paths
	perms = fs.FileMode(0777)
)

func CreateObject(logger logger.Logger, bucketName string, objectKey string) (*Object, error) {
	return &Object{
		logger:     logger,
		ObjectKey:  objectKey,
		BucketName: bucketName,
		Perms:      perms,
	}, nil
}

func (o *Object) SetCachePath(cachePath string) {
	o.Path = cachePath + "/" + o.BucketName + "/" + o.ObjectKey
}

func (o *Object) SaveByteStreamToFile(objectStream []byte) error {
	// Check if the parent directory of the file exists, and create it if it doesn't exist
	o.logger.Debugw("Cache path",
		"cachePath", o.Path,
	)
	if o.Path == "" {
		o.logger.Errorw("Cache path is empty",
			"cachePath", o.Path,
		)
		return errors.New("Error cache path is empty")
	}
	destDir := filepath.Dir(o.Path)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, o.Perms); err != nil {
			o.logger.Errorw("Failed to create parent directory",
				"error", err,
			)
			return err
		}
	}

	// Write the byte stream to the file
	err := ioutil.WriteFile(o.Path, objectStream, o.Perms)
	if err != nil {
		o.logger.Errorw("Failed to save byte stream to file",
			"error", err,
		)
		return err
	}

	return nil
}

func (o *Object) RemoveFileFromCache() error {
	o.logger.Debugw("Removing file from cache",
		"cachePath", o.Path,
	)
	err := os.Remove(o.Path)
	if err != nil {
		o.logger.Errorw("Failed to remove file from cache",
			"error", err,
		)
		return err
	}
	return nil
}
