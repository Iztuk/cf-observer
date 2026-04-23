package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func LoadConfigFile(override string) (map[string]Host, error) {
	var configPath string

	if override != "" {
		configPath = override
	} else {
		baseDir, err := os.UserConfigDir()
		if err != nil {
			return map[string]Host{}, err
		}

		configPath = filepath.Join(baseDir, "codeforge-observer", "config.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return map[string]Host{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return map[string]Host{}, err
	}

	err = config.Validate()
	if err != nil {
		return map[string]Host{}, err
	}

	AppRunTimeConfig = config.RunTime

	return config.ValidateHostUrls()
}

func (c *Config) Validate() error {
	if c.RunTime.Listen == "" {
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

	if c.RunTime.AuditConfig.QueueSize < 0 {
		return fmt.Errorf("audit queue_size cannot be negative")
	}
	if c.RunTime.AuditConfig.Workers < 0 {
		return fmt.Errorf("audit workers cannot be negative")
	}

	return nil
}

func (c *Config) ValidateHostUrls() (map[string]Host, error) {
	for key, host := range c.Hosts {
		u, err := url.Parse(host.UpstreamRaw)
		if err != nil {
			return map[string]Host{}, err
		}

		c.Hosts[key] = Host{
			UpstreamRaw:      host.UpstreamRaw,
			Upstream:         u,
			ApiContract:      host.ApiContract,
			ResourceContract: host.ResourceContract,
		}
	}

	return c.Hosts, nil
}
