package daemon

import (
	"cf-observer/internal/config"
	"cf-observer/internal/proxy"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func RunDaemon(hosts map[string]config.Host) error {
	f, err := os.OpenFile(config.AppProcessConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	logger := log.New(f, "cf-observer: ", log.LstdFlags)

	pm, err := proxy.NewProxyManager(hosts, logger)
	if err != nil {
		return fmt.Errorf("create proxy manager: %w", err)
	}

	server := &http.Server{
		Addr:    config.AppRunTimeConfig.Listen,
		Handler: pm,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %v", err)
	}

	return nil
}
