package scanner

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type scanCollector struct {
	sugar         *zap.SugaredLogger
	filesScanned  prometheus.Counter
	infectedFiles prometheus.Counter
	cleanFiles    prometheus.Counter
	scanErrors    prometheus.Counter
	scanTime      prometheus.Histogram
}

func CreateScanCollector(sugar *zap.SugaredLogger) (*scanCollector, error) {
	filesScanned := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scanner_total_scans", Help: "Total number of scans performed by Aegis"})
	cleanFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scanner_infected_files", Help: "Total of infected files scanned by Aegis"})
	infectedFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scanner_clean_files", Help: "Total of clean files scanned by Aegis"})
	scanErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scanner_errors", Help: "Total number of errors encountered during scans by Aegis"})
	scanTime := promauto.NewHistogram(prometheus.HistogramOpts{Name: "aegis_scanner_time", Help: "Time taken to perform a scan", Buckets: []float64{0, 125, 250, 500, 1000, 2000, 4000, 8000, 16000}})
	return &scanCollector{
		sugar:         sugar,
		filesScanned:  filesScanned,
		infectedFiles: infectedFiles,
		cleanFiles:    cleanFiles,
		scanErrors:    scanErrors,
		scanTime:      scanTime,
	}, nil
}

func (c *scanCollector) FileScanned() {
	c.sugar.Debug("Incrementing files scanned counter")
	c.filesScanned.Inc()
}

func (c *scanCollector) CleanFile() {
	c.sugar.Debug("Incrementing clean files scanned counter")
	c.cleanFiles.Inc()
}

func (c *scanCollector) InfectedFile() {
	c.sugar.Debug("Incrementing infected files scanned counter")
	c.infectedFiles.Inc()
}

func (c *scanCollector) ScanError() {
	c.sugar.Debug("Incrementing scan errors counter")
	c.scanErrors.Inc()
}

func (c *scanCollector) ScanTime(t float64) {
	c.sugar.Debug("Incrementing scan time histogram")
	c.sugar.Debugw("Scan time", "time", t)
	c.scanTime.Observe(t)
}
