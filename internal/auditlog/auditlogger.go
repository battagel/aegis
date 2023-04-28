package auditlog

import (
	"aegis/pkg/logger"
)

type Database interface {
	CreateTable(string) error
	Insert(string, string, string, string, string, string, string) error
}

type AuditLogger struct {
	logger    logger.Logger
	db        Database
	tableName string
}

func CreateAuditLogger(logger logger.Logger, db Database, tableName string) (*AuditLogger, error) {
	// Check that db and table exists
	db.CreateTable(tableName)
	return &AuditLogger{
		logger:    logger,
		db:        db,
		tableName: tableName,
	}, nil
}

func (a *AuditLogger) Log(bucketName, objectKey, result, antivirus, timestamp, virusType string) {
	a.logger.Debugln("Adding audit log")
	a.db.Insert(a.tableName, bucketName, objectKey, result, antivirus, timestamp, virusType)
}
