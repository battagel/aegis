package objectstore

import (
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
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
	mockMinio   *MockMinio
	objectStore *ObjectStore
}

func CreateTestItems(t *testing.T) *testItems {
	logger, err := logger.CreateZapLogger("debug", "console")

	mockMinio := new(MockMinio)
	mockObjectStoreCollector := new(MockObjectStoreCollector)
	mockObjectStoreCollector.On("GetObject").Return(nil)
	mockObjectStoreCollector.On("PutObject").Return(nil)
	mockObjectStoreCollector.On("RemoveObject").Return(nil)
	mockObjectStoreCollector.On("GetObjectTagging").Return(nil)
	mockObjectStoreCollector.On("PutObjectTagging").Return(nil)

	objectStore, err := CreateObjectStore(logger, mockMinio, mockObjectStoreCollector)
	assert.Nil(t, err)

	return &testItems{
		logger:      logger,
		mockMinio:   mockMinio,
		objectStore: objectStore,
	}
}

func TestObjectStore_GetObject_Happy(t *testing.T) {
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("GetObject").Return(nil)

	for _, test := range tableTests {
		testItems.mockMinio.On("GetObject", test.bucketName, test.objectKey).Return([]byte("test"), nil)
		_, err := testItems.objectStore.GetObject(test.bucketName, test.objectKey)
		assert.Nil(t, err)
	}
}

func TestObjectStore_PutObject_Happy(t *testing.T) {
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("PutObject").Return(nil)

	for _, test := range tableTests {
		testItems.mockMinio.On("PutObject", test.bucketName, test.objectKey, mock.AnythingOfType("[]uint8")).Return(nil)
		err := testItems.objectStore.PutObject(test.bucketName, test.objectKey, []byte("test"))
		assert.Nil(t, err)
	}
}

func TestObjectStore_RemoveObject_Happy(t *testing.T) {
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("RemoveObject").Return(nil)

	for _, test := range tableTests {
		testItems.mockMinio.On("RemoveObject", test.bucketName, test.objectKey).Return(nil)
		err := testItems.objectStore.RemoveObject(test.bucketName, test.objectKey)
		assert.Nil(t, err)
	}
}

func TestObjectStore_GetObjectTagging_Happy(t *testing.T) {
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("GetObjectTagging").Return(nil)

	for _, test := range tableTests {
		testItems.mockMinio.On("GetObjectTagging", test.bucketName, test.objectKey).Return(map[string]string{"test": "test"}, nil)
		_, err := testItems.objectStore.GetObjectTagging(test.bucketName, test.objectKey)
		assert.Nil(t, err)
	}
}

func TestObjectStore_PutObjectTagging_Happy(t *testing.T) {
	testItems := CreateTestItems(t)

	testItems.mockMinio.On("PutObjectTagging").Return(nil)

	for _, test := range tableTests {
		testItems.mockMinio.On("PutObjectTagging", test.bucketName, test.objectKey, mock.AnythingOfType("map[string]string")).Return(nil)
		err := testItems.objectStore.PutObjectTagging(test.bucketName, test.objectKey, map[string]string{"test": "test"})
		assert.Nil(t, err)
	}
}
