package daemon

import (
	"cf-observer/internal/config"
	"cf-observer/internal/proxy"
	"fmt"
	"log"
	"os"
)

func RunDaemon(hosts map[string]config.Host) error {
	f, err := os.OpenFile(config.AppProcessConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	logger := log.New(f, "cf-observer: ", log.LstdFlags)

	// Ensure single instance of daemon

	_, err = proxy.NewProxyManager(hosts, logger)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}
