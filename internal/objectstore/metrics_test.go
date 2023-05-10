package objectstore

import (
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanCollector_Happy(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	objectStoreCollector, err := CreateObjectStoreCollector(logger)
	assert.Nil(t, err)

	objectStoreCollector.GetObject()
	objectStoreCollector.PutObject()
	objectStoreCollector.RemoveObject()
	objectStoreCollector.GetObjectTagging()
	objectStoreCollector.PutObjectTagging()
}
