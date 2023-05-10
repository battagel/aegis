package objectstore

import (
	"aegis/mocks"
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tableTest struct {
	bucketName, objectKey string
}

var tableTests = []tableTest{
	{"test", "test"},
	{"1234", "1234"},
	{"test bucket", "test object"},
	{"test-bucket", "test-object"},
	{"test_bucket", "test_object"},
}

type testItems struct {
	logger      logger.Logger
	objectStore *ObjectStore
}

func CreateTestItems(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")

	mockMinio := new(mocks.Minio)
	mockObjectStoreCollector := new(mocks.ObjectStoreCollector)
	mockObjectStoreCollector.On("GetObject").Return(nil)
	mockObjectStoreCollector.On("PutObject").Return(nil)
	mockObjectStoreCollector.On("RemoveObject").Return(nil)
	mockObjectStoreCollector.On("GetObjectTagging").Return(nil)
	mockObjectStoreCollector.On("PutObjectTagging").Return(nil)

	objectStore, err := CreateObjectStore(logger, mockMinio, mockObjectStoreCollector)
	assert.Nil(t, err)

	return &testItems{
		logger:      logger,
		objectStore: objectStore,
	}
}

func TestObjectStore_GetObject_Happy()
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("GetObject", ).Return(nil)

	testItems.objectStore.GetObject()

	
}


