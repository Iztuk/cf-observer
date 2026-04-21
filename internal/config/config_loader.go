package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func LoadConfigFile(override string) (*Config, error) {
	var configPath string

	if override != "" {
		configPath = override
	} else {
		baseDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		configPath = filepath.Join(baseDir, "codeforge-observer", "config.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Listen == "" {
		return fmt.Errorf("listen is required")
	}

	for hostName, host := range c.Hosts {
		if hostName == "" {
			return fmt.Errorf("host name cannot be empty")
		}
		if host.Upstream == nil {
			return fmt.Errorf("host %q is missing upstream", hostName)
		}
	}

	if c.AuditConfig.QueueSize < 0 {
		return fmt.Errorf("audit queue_size cannot be negative")
	}
	if c.AuditConfig.Workers < 0 {
		return fmt.Errorf("audit workers cannot be negative")
	}

	return nil
}
