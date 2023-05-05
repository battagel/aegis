package cleaner

import (
	"aegis/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type cleanerCollector struct {
	logger             logger.Logger
	objectsRemoved     prometheus.Counter
	objectsTagged      prometheus.Counter
	objectsQuarentined prometheus.Counter
	cleanupErrors      prometheus.Counter
}

func CreateCleanerCollector(logger logger.Logger) (*cleanerCollector, error) {
	objectsRemoved := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_cleaner_object_removed", Help: "Cleaner total objects removed"})
	objectsTagged := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_cleaner_object_tagged", Help: "Cleaner total objects tagging"})
	objectsQuarentined := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_cleaner_object_quarentined", Help: "Cleaner total objects quarentined"})
	cleanupErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_cleaner_errors", Help: "Cleaner total errors"})
	return &cleanerCollector{
		logger:             logger,
		objectsRemoved:     objectsRemoved,
		objectsTagged:      objectsTagged,
		objectsQuarentined: objectsQuarentined,
		cleanupErrors:      cleanupErrors,
	}, nil

}

// Metric update functions
func (c *cleanerCollector) ObjectRemoved() {
	c.logger.Debugw("Incrementing cleaner object removed counter")
	c.objectsRemoved.Inc()
}

func (c *cleanerCollector) ObjectTagged() {
	c.logger.Debugw("Incrementing cleaner object tagged counter")
	c.objectsTagged.Inc()
}

func (c *cleanerCollector) ObjectQuarantined() {
	c.logger.Debugw("Incrementing cleaner object quarentined counter")
	c.objectsQuarentined.Inc()
}

func (c *cleanerCollector) CleanupError() {
	c.logger.Debugw("Incrementing cleaner error counter")
	c.cleanupError.Inc()
}
