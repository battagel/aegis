package clamav

import (
	"aegis/pkg/logger"
	"os/exec"
	"regexp"
	"strings"
)

type ClamAVScanner struct {
	logger logger.Logger
}

func CreateClamAV(logger logger.Logger) (*ClamAVScanner, error) {
	return &ClamAVScanner{
		logger: logger,
	}, nil
}

func (c *ClamAVScanner) ScanFile(filePath string) (bool, string, error) {
	// Returns false if file is clean, true if infected
	// If there are any errors then return true (infected)
	clamCmd := exec.Command("clamdscan", filePath, "--config=clamd.conf")
	output, err := clamCmd.Output()
	c.logger.Debugw("clamdscan output",
		"output", string(output),
	)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means infected
			if exitError.ExitCode() == 1 {
				virusType := c.findVirusType(string(output))
				return true, virusType, nil
			}
		}
		c.logger.Errorw("Error running clamdscan. Is clamd running?",
			"error", err,
		)
		return true, "", err
	}
	// Due to exit codes, the file must be ok
	return false, "", nil
}

func (c *ClamAVScanner) findVirusType(output string) string {
	re := regexp.MustCompile(`(\w+.\w+.\w+-\w+-\w+\w+-\w+) FOUND`)
	virusType := strings.TrimSuffix(re.FindString(output), " FOUND")
	c.logger.Debugw("virus type output",
		"virusType", virusType,
	)
	return virusType
}

func (c *ClamAVScanner) GetName() string {
	return "clamav"
}
