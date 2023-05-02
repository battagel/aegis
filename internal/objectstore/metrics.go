package objectstore

import (
	"aegis/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type objectStoreCollector struct {
	logger            logger.Logger
	getObjects        prometheus.Counter
	putObjects        prometheus.Counter
	removeObjects     prometheus.Counter
	getObjectsTagging prometheus.Counter
	putObjectsTagging prometheus.Counter
}

func CreateObjectStoreCollector(logger logger.Logger) (*objectStoreCollector, error) {
	getObjects := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_get_objects", Help: "Object Store total get objects"})
	putObjects := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_put_objects", Help: "Object Store total put objects"})
	removeObjects := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_remove_objects", Help: "Object Store total removed objects"})
	getObjectsTagging := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_get_objects_tagging", Help: "Object Store total get objects tagging"})
	putObjectsTagging := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_put_objects_tagging", Help: "Object Store total put objects tagging"})
	return &objectStoreCollector{
		logger:            logger,
		getObjects:        getObjects,
		putObjects:        putObjects,
		removeObjects:     removeObjects,
		getObjectsTagging: getObjectsTagging,
		putObjectsTagging: putObjectsTagging,
	}, nil
}

// Metric update functions
func (c *objectStoreCollector) GetObject() {
	c.logger.Debugln("Incrementing get object counter")
	c.getObjects.Inc()
}

func (c *objectStoreCollector) PutObject() {
	c.logger.Debugln("Incrementing put object counter")
	c.putObjects.Inc()
}

func (c *objectStoreCollector) RemoveObject() {
	c.logger.Debugln("Incrementing remove object counter")
	c.removeObjects.Inc()
}

func (c *objectStoreCollector) GetObjectTagging() {
	c.logger.Debugln("Incrementing get object tagging counter")
	c.getObjectsTagging.Inc()
}

func (c *objectStoreCollector) PutObjectTagging() {
	c.logger.Debugln("Incrementing put object tagging counter")
	c.putObjectsTagging.Inc()
}
