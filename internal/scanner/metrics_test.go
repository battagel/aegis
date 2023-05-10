package scanner

import (
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanCollector_Happy(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	scanCollector, err := CreateScanCollector(logger)
	assert.Nil(t, err)

	scanCollector.FileScanned()
	scanCollector.CleanFile()
	scanCollector.InfectedFile()
	scanCollector.ScanError()
	scanCollector.ScanTime(0.1)
}
