package cleaner

import (
	"aegis/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanCollector_Happy(t *testing.T) {
	logger, err := logger.CreateZapLogger("debug", "console")
	assert.Nil(t, err)

	cleanerCollector, err := CreateCleanerCollector(logger)
	assert.Nil(t, err)

	cleanerCollector.ObjectRemoved()
	cleanerCollector.ObjectTagged()
	cleanerCollector.ObjectQuarantined()
	cleanerCollector.CleanupError()
}
