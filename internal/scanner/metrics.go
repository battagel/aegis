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
}

func CreateScanCollector(sugar *zap.SugaredLogger) (*scanCollector, error) {
	filesScanned := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_total_scans", Help: "Total number of scans performed by Aegis"})
	cleanFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_infected_files_scanned", Help: "Total of infected files scanned by Aegis"})
	infectedFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_clean_files_scanned", Help: "Total of clean files scanned by Aegis"})
	scanErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scan_errors", Help: "Total number of errors encountered during scans by Aegis"})
	return &scanCollector{
		sugar:         sugar,
		filesScanned:  filesScanned,
		infectedFiles: infectedFiles,
		cleanFiles:    cleanFiles,
		scanErrors:    scanErrors,
	}, nil
}

// Covered by promautos in-built registry
// func (c *scanCollector) GatherCollector() {
//
// }
//
// func (c *scanCollector) DescribeCollector() {
//
// }

// Metric update functions
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
