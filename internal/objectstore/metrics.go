package objectstore

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type objectStoreCollector struct {
	sugar             *zap.SugaredLogger
	getObjects        prometheus.Counter
	getObjectsTagging prometheus.Counter
	putObjectsTagging prometheus.Counter
}

func CreateObjectStoreCollector(sugar *zap.SugaredLogger) (*objectStoreCollector, error) {
	getObjects := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_get_objects", Help: "Object Store total get objects"})
	getObjectsTagging := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_get_objects_tagging", Help: "Object Store total get objects tagging"})
	putObjectsTagging := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_objectstore_put_objects_tagging", Help: "Object Store total put objects tagging"})
	return &objectStoreCollector{
		sugar:             sugar,
		getObjects:        getObjects,
		getObjectsTagging: getObjectsTagging,
		putObjectsTagging: putObjectsTagging,
	}, nil
}

// Metric update functions
func (c *objectStoreCollector) GetObject() {
	c.sugar.Debugln("Incrementing get object counter")
	c.getObjects.Inc()
}

func (c *objectStoreCollector) GetObjectTagging() {
	c.sugar.Debugln("Incrementing get object tagging counter")
	c.getObjectsTagging.Inc()
}

func (c *objectStoreCollector) PutObjectTagging() {
	c.sugar.Debugln("Incrementing put object tagging counter")
	c.putObjectsTagging.Inc()
}
