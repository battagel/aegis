package auditlog

import (
	"aegis/mocks"
	"aegis/pkg/logger"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuditLogger_Log_Happy(t *testing.T) {

	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	tableName := "test"
	mockDB := new(mocks.Database)
	mockDB.On("CreateTable", tableName).Return(nil)
	mockDB.On("Insert", tableName, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	auditLogger, err := CreateAuditLogger(logger, mockDB, tableName)
	assert.Nil(t, err)

	auditLogger.Log("test", "test", "test", "test", "test", "test")
}

func TestAuditLogger_Log_Error(t *testing.T) {

	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	tableName := "test"
	mockDB := new(mocks.Database)
	mockDB.On("CreateTable", tableName).Return(nil)
	mockDB.On("Insert", tableName, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Error inserting into table"))

	auditLogger, err := CreateAuditLogger(logger, mockDB, tableName)
	assert.Nil(t, err)

	auditLogger.Log("test", "test", "test", "test", "test", "test")
}
