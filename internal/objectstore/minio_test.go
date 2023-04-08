package objectstore

import (
	"aegis/mocks"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/tags"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type CommonTestItems struct {
	Sugar *zap.SugaredLogger
}

func ProvideCommonTestItems(t *testing.T) *CommonTestItems {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Error creating logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return &CommonTestItems{Sugar: sugar}
}

func TestObjectStore(t *testing.T) {
	commonTestItems := ProvideCommonTestItems(t)
	mockMinio := new(mocks.MinioClient)
	mocksObjectStoreCollector := new(mocks.ObjectStoreCollector)

	objectStore, err := CreateObjectStore(commonTestItems.Sugar, mockMinio, mocksObjectStoreCollector)
	assert.NoError(t, err)

	// GetObject
	goodObject := &minio.Object{}
	goodData, err := io.ReadAll(goodObject)

	mockMinio.On("GetObject", context.Background(), "test-bucket", "good.txt", minio.GetObjectOptions{}).Return(*goodObject, nil)
	object, err := objectStore.GetObject("test-bucket", "good.txt")
	assert.NoError(t, err)
	assert.Equal(t, goodData, object)

	err = errors.New("No object")
	mockMinio.On("GetObject", context.Background(), "test-bucket", "", minio.GetObjectOptions{}).Return(nil, err)
	object, err = objectStore.GetObject("test-bucket", "")
	assert.Error(t, err)
	assert.Equal(t, "", object)

	err = errors.New("No bucket")
	mockMinio.On("GetObject", context.Background(), "", "text.txt", minio.GetObjectOptions{}).Return(nil, err)
	object, err = objectStore.GetObject("", "good.txt")
	assert.Error(t, err)
	assert.Equal(t, "", object)

	// GetObjectTagging
	// goodTags, err := tags.MapToObjectTags(map[string]string{"test": "test"})
	// assert.NoError(t, err)

	// mockMinio.On("GetObjectTagging", "test-bucket", "good.txt").Return(goodTags, nil)
	// tags, err := objectStore.GetObjectTagging("test-bucket", "good.txt")
	// assert.NoError(t, err)
	// assert.Equal(t, goodTags, tags)

	// err = errors.New("No object")
	// mockMinio.On("GetObjectTagging", "test-bucket", "").Return(nil, err)
	// tags, err = objectStore.GetObjectTagging("test-bucket", "good.txt")
	// assert.Error(t, err)
	// assert.Equal(t, goodTags, tags)

	// err = errors.New("No bucket")
	// mockMinio.On("GetObjectTagging", "", "good.txt").Return(nil, err)
	// tags, err = objectStore.GetObjectTagging("test-bucket", "good.txt")
	// assert.Error(t, err)
	// assert.Equal(t, goodTags, tags)

	// // PutObjectTagging
	// mockMinio.On("PutObjectTagging", "test-bucket", "good.txt").Return(nil)
	// err = objectStore.PutObjectTagging("test-bucket", "good.txt", goodTags)
	// assert.NoError(t, err)

	// err = errors.New("No object")
	// mockMinio.On("PutObjectTagging", "test-bucket", "").Return(err)
	// err = objectStore.PutObjectTagging("test-bucket", "good.txt", goodTags)
	// assert.Error(t, err)

	// err = errors.New("No bucket")
	// mockMinio.On("PutObjectTagging", "", "good.txt").Return(err)
	// err = objectStore.PutObjectTagging("test-bucket", "good.txt", goodTags)
	// assert.Error(t, err)
}
