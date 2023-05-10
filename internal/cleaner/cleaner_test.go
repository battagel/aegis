package cleaner

import (
	"aegis/internal/object"
	"aegis/mocks"
	"aegis/pkg/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type tableTest struct {
	bucketName, objectKey string
	result                bool
}

var tableTests = []tableTest{
	{"test", "test", false},
	{"test", "test", true},
	// {"1234", "1234", },
	// {"test bucket", "test object"},
	// {"test-bucket", "test-object"},
	// {"test_bucket", "test_object"},
}

func TestEventsManager_Start_Remove_Happy(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	mockObjectStore := new(mocks.ObjectStore)
	mockObjectStore.On("RemoveObject", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	cleanupPolicy := "remove"
	quarantineBucket := ""
	mockCleanerCollector := new(mocks.CleanerCollector)
	mockCleanerCollector.On("ObjectRemoved").Return(nil)
	mockCleanerCollector.On("ObjectTagged").Return(nil)
	mockCleanerCollector.On("ObjectQuarantined").Return(nil)
	mockCleanerCollector.On("CleanupError").Return(nil)

	mockAuditLogger := new(mocks.AuditLogger)
	mockAuditLogger.On("Log", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	cleaner, err := CreateCleaner(logger, mockObjectStore, cleanupPolicy, quarantineBucket, mockCleanerCollector, mockAuditLogger)

	scanTime := time.Now().Format("01-02-2006 15:04:05")

	for _, test := range tableTests {
		request, err := object.CreateObject(logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		err = cleaner.Cleanup(request, test.result, scanTime)
		assert.Nil(t, err)
	}
}

func TestEventsManager_Start_Tag_Happy(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	mockObjectStore := new(mocks.ObjectStore)
	mockObjectStore.On("AddObjectTagging", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(nil)
	cleanupPolicy := "tag"
	quarantineBucket := ""
	mockCleanerCollector := new(mocks.CleanerCollector)
	mockCleanerCollector.On("ObjectRemoved").Return(nil)
	mockCleanerCollector.On("ObjectTagged").Return(nil)
	mockCleanerCollector.On("ObjectQuarantined").Return(nil)
	mockCleanerCollector.On("CleanupError").Return(nil)

	mockAuditLogger := new(mocks.AuditLogger)
	mockAuditLogger.On("Log", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	cleaner, err := CreateCleaner(logger, mockObjectStore, cleanupPolicy, quarantineBucket, mockCleanerCollector, mockAuditLogger)

	scanTime := time.Now().Format("01-02-2006 15:04:05")

	for _, test := range tableTests {
		request, err := object.CreateObject(logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		err = cleaner.Cleanup(request, test.result, scanTime)
		assert.Nil(t, err)
	}
}

func TestEventsManager_Start_Quarantine_Happy(t *testing.T) {
	// Create mock logger
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	mockObjectStore := new(mocks.ObjectStore)
	mockObjectStore.On("MoveObject", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	cleanupPolicy := "quarantine"
	quarantineBucket := "quarantine"
	mockCleanerCollector := new(mocks.CleanerCollector)
	mockCleanerCollector.On("ObjectRemoved").Return(nil)
	mockCleanerCollector.On("ObjectTagged").Return(nil)
	mockCleanerCollector.On("ObjectQuarantined").Return(nil)
	mockCleanerCollector.On("CleanupError").Return(nil)

	mockAuditLogger := new(mocks.AuditLogger)
	mockAuditLogger.On("Log", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	cleaner, err := CreateCleaner(logger, mockObjectStore, cleanupPolicy, quarantineBucket, mockCleanerCollector, mockAuditLogger)

	scanTime := time.Now().Format("01-02-2006 15:04:05")

	for _, test := range tableTests {
		request, err := object.CreateObject(logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		err = cleaner.Cleanup(request, test.result, scanTime)
		assert.Nil(t, err)
	}
}
