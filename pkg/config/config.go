package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Vault VaultConfig `yaml:"vault"`
	K8s   K8sConfig   `yaml:"kubernetes"`
	App   AppConfig   `yaml:"app"`
}

// VaultConfig contains Vault-specific configuration
type VaultConfig struct {
	Nodes     []string  `yaml:"nodes,omitempty"`
	Namespace string    `yaml:"namespace,omitempty"`
	Keys      []string  `yaml:"unseal_keys"`
	TLS       TLSConfig `yaml:"tls"`
}

// TLSConfig contains TLS configuration for Vault
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CACert     string `yaml:"ca_cert,omitempty"`
	ClientCert string `yaml:"client_cert,omitempty"`
	ClientKey  string `yaml:"client_key,omitempty"`
	Insecure   bool   `yaml:"insecure,omitempty"`
}

// K8sConfig contains Kubernetes-specific configuration
type K8sConfig struct {
	InCluster  bool   `yaml:"in_cluster"`
	ConfigPath string `yaml:"config_path,omitempty"`
	Namespace  string `yaml:"namespace"`
}

// AppConfig contains general application configuration
type AppConfig struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	LogLevel      string        `yaml:"log_level"`
	MetricsAddr   string        `yaml:"metrics_addr,omitempty"`
}

// Load reads and parses the configuration file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.App.CheckInterval == 0 {
		config.App.CheckInterval = 30 * time.Second
	}
	if config.App.LogLevel == "" {
		config.App.LogLevel = "info"
	}
	if config.K8s.Namespace == "" {
		config.K8s.Namespace = "vault"
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Vault.Nodes) == 0 {
		return fmt.Errorf("vault address is required")
	}
	if len(c.Vault.Keys) == 0 {
		return fmt.Errorf("at least one unseal key is required")
	}
	return nil
}
