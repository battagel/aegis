package scanner

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type scanCollector struct {
	filesScanned  prometheus.Counter
	infectedFiles prometheus.Counter
	cleanFiles    prometheus.Counter
	scanErrors    prometheus.Counter
}

func CreateScanCollector() (*scanCollector, error) {
	filesScanned := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_total_scans", Help: "Total number of scans performed by Aegis"})
	cleanFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_infected_files_scanned", Help: "Total of infected files scanned by Aegis"})
	infectedFiles := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_clean_files_scanned", Help: "Total of clean files scanned by Aegis"})
	scanErrors := promauto.NewCounter(prometheus.CounterOpts{Name: "aegis_scan_errors", Help: "Total number of errors encountered during scans by Aegis"})
	return &scanCollector{
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
	c.filesScanned.Inc()
}

func (c *scanCollector) CleanFile() {
	c.cleanFiles.Inc()
}

func (c *scanCollector) InfectedFile() {
	c.infectedFiles.Inc()
}

func (c *scanCollector) ScanError() {
	c.scanErrors.Inc()
}
