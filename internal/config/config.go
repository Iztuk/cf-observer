package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type Config struct {
	Listen      string          `yaml:"listen"`
	DefaultHost *Host           `yaml:"default_host"`
	Hosts       map[string]Host `yaml:"hosts"`
	AuditConfig AuditConfig     `yaml:"audit_config"`
}

type Host struct {
	UpstreamRaw      string   `yaml:"upstream"`
	Upstream         *url.URL `yaml:"-"`
	ApiContract      string   `yaml:"api_contract"`
	ResourceContract string   `yaml:"resource_contract"`
}

type AuditConfig struct {
	Enabled   bool `yaml:"enabled"`
	QueueSize int  `yaml:"queue_size"`
	Workers   int  `yaml:"worker"`
}

var AppConfig *Config

const defaultConfigYAML = `listen: ":8080"

default_host:
  upstream: "http://localhost:8081"
  api_contract: "contracts/default.json"
  resource_contract: "resources/default.json"

hosts: {}

audit_config:
  enabled: true
  queue_size: 1000
  workers: 4
`

func ConfigDir() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}

	return filepath.Join(baseDir, "codeforge-observer"), nil
}

func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.yaml"), nil
}

func InitConfigDir(force bool) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create config directory: %w", err)
	}

	configPath := filepath.Join(dir, "config.yaml")

	if !force {
		_, err = os.Stat(configPath)
		if err == nil {
			return configPath, fmt.Errorf("config file already exists: %s", configPath)
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("check config file: %w", err)
		}
	}

	if err := os.WriteFile(configPath, []byte(defaultConfigYAML), 0o644); err != nil {
		return "", fmt.Errorf("write config file: %w", err)
	}

	return configPath, nil
}
