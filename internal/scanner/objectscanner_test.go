package scanner

import (
	"aegis/internal/object"
	"aegis/mocks"
	"aegis/pkg/logger"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	logger            logger.Logger
	mockObjectStore   *mocks.ObjectStore
	mockAntivirus     *mocks.Antivirus
	mockCleaner       *mocks.Cleaner
	mockAuditLogger   *mocks.AuditLogger
	mockScanCollector *mocks.ScanCollector
	scanner           *Scanner
	errChan           chan error
}

func createTestItems(t *testing.T) *testItems {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	removeAfterScan := true
	datetimeFormat := "01-02-2006 15:04:05"
	cachePath := "cache/testing"

	mockObjectStore := new(mocks.ObjectStore)
	mockAntivirus := new(mocks.Antivirus)
	mockCleaner := new(mocks.Cleaner)
	mockAuditLogger := new(mocks.AuditLogger)
	mockScanCollector := new(mocks.ScanCollector)

	errChan := make(chan error, 1)

	scanner, err := CreateObjectScanner(logger, mockObjectStore, []Antivirus{mockAntivirus}, mockCleaner, mockAuditLogger, mockScanCollector, removeAfterScan, datetimeFormat, cachePath)
	assert.Nil(t, err)

	mockAuditLogger.On("Log", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	mockScanCollector.On("FileScanned").Return(nil)
	mockScanCollector.On("CleanFile").Return(nil)
	mockScanCollector.On("InfectedFile").Return(nil)
	mockScanCollector.On("ScanError").Return(nil)
	mockAntivirus.On("GetName").Return("mock-av")
	mockCleaner.On("Cleanup", mock.AnythingOfType("*object.Object"), mock.AnythingOfType("bool"), mock.AnythingOfType("string")).Return(nil)

	return &testItems{
		logger:            logger,
		mockObjectStore:   mockObjectStore,
		mockAntivirus:     mockAntivirus,
		mockCleaner:       mockCleaner,
		mockAuditLogger:   mockAuditLogger,
		mockScanCollector: mockScanCollector,
		scanner:           scanner,
		errChan:           errChan,
	}
}

func TestObjectScanner_ScanObject_Happy(t *testing.T) {
	testItems := createTestItems(t)

	for _, test := range tableTests {
		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte("test"), nil)
		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(false, "", nil)

		testItems.scanner.ScanObject(request, testItems.errChan)
		select {
		case err := <-testItems.errChan:
			testItems.logger.Errorw("Test failed",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
			t.Fail()
		default:
			testItems.logger.Infow("Test passed",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
			)
		}
	}
	cleanupCache()
}

func TestObjectScanner_ScanObject_Infected(t *testing.T) {
	testItems := createTestItems(t)

	for _, test := range tableTests {
		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte("test"), nil)
		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(true, "Win-virus", nil)

		testItems.scanner.ScanObject(request, testItems.errChan)
		select {
		case err := <-testItems.errChan:
			testItems.logger.Errorw("Test failed",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
			t.Fail()
		default:
			testItems.logger.Infow("Test passed",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
			)
		}
	}
	cleanupCache()
}

// Errors
func TestObjectScanner_ScanObject_GetObjectError(t *testing.T) {
	testItems := createTestItems(t)

	for _, test := range tableTests {
		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte(""), errors.New("Error Getting Object"))
		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(true, "Win-virus", nil)

		testItems.scanner.ScanObject(request, testItems.errChan)
		select {
		case err := <-testItems.errChan:
			assert.EqualError(t, err, "Error Getting Object")
			testItems.logger.Errorw("Test passed: Error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
		default:
			testItems.logger.Infow("Test failed: No error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
			)
			t.Fail()
		}

	}
	cleanupCache()
}

// func TestObjectScanner_ScanObject_SaveFileError(t *testing.T) {
// 	testItems := createTestItems(t)

// 	for _, test := range tableTests {
// 		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
// 		assert.Nil(t, err)

// 		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte(""), nil)
// 		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(false, "", nil)

// 		testItems.scanner.ScanObject(request, testItems.errChan)
// 		select {
// 		case err := <-testItems.errChan:
// 			assert.EqualError(t, err, "remove "+testItems.scanner.cachePath+"/"+test.bucketName+"/"+test.objectKey+": no such file or directory")
// 			testItems.logger.Errorw("Test passed: Error returned",
// 				"bucketName", test.bucketName,
// 				"objectKey", test.objectKey,
// 				"err", err,
// 			)
// 		default:
// 			testItems.logger.Infow("Test failed: No error returned",
// 				"bucketName", test.bucketName,
// 				"objectKey", test.objectKey,
// 			)
// 			t.Fail()
// 		}
// 	}
//  cleanupCache()
// }

func TestObjectScanner_ScanObject_ScanFileError(t *testing.T) {
	testItems := createTestItems(t)

	for _, test := range tableTests {
		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte(""), nil)
		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(true, "", errors.New("Error Scanning File"))

		testItems.scanner.ScanObject(request, testItems.errChan)
		select {
		case err := <-testItems.errChan:
			assert.EqualError(t, err, "Error Scanning File")
			testItems.logger.Errorw("Test passed: Error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
		default:
			testItems.logger.Infow("Test failed: No error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
			)
			t.Fail()
		}
	}
	cleanupCache()
}

func TestObjectScanner_ScanObject_RemoveFileFromCache(t *testing.T) {
	testItems := createTestItems(t)

	for _, test := range tableTests {
		request, err := object.CreateObject(testItems.logger, test.bucketName, test.objectKey)
		assert.Nil(t, err)

		testItems.mockObjectStore.On("GetObject", test.bucketName, test.objectKey).Return([]byte(""), nil)
		testItems.mockAntivirus.On("ScanFile", mock.AnythingOfType("string")).Return(false, "", nil)

		go func() {
			for {
				os.Remove(testItems.scanner.cachePath + "/" + test.bucketName + "/" + test.objectKey)
			}
		}()

		testItems.scanner.ScanObject(request, testItems.errChan)
		select {
		case err := <-testItems.errChan:
			assert.EqualError(t, err, "remove "+testItems.scanner.cachePath+"/"+test.bucketName+"/"+test.objectKey+": no such file or directory")
			testItems.logger.Errorw("Test passed: Error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
				"err", err,
			)
		default:
			testItems.logger.Infow("Test failed: No error returned",
				"bucketName", test.bucketName,
				"objectKey", test.objectKey,
			)
			t.Fail()
		}
	}
	cleanupCache()
}

func cleanupCache() {
	os.RemoveAll("cache")
}
