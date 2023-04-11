package clamav

import (
	"aegis/pkg/logger"
	"os/exec"
)

type ClamAVScanner struct {
	logger logger.Logger
}

func CreateClamAV(logger logger.Logger) (*ClamAVScanner, error) {
	return &ClamAVScanner{
		logger: logger,
	}, nil
}

func (c *ClamAVScanner) ScanFile(filePath string) (bool, error) {
	// Returns false if file is clean, true if infected
	// If there are any errors then return true (infected)
	cmd := exec.Command("clamdscan", filePath, "--stream", "-m")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means infected
			if exitError.ExitCode() == 1 {
				return true, nil
			}
		}
		c.logger.Errorw("Error running clamdscan. Is clamd running?",
			"error", err,
		)
		return true, err
	}
	// Due to exit codes, the file must be ok
	return false, nil
}

func (c *ClamAVScanner) GetName() string {
	return "clamav"
}
