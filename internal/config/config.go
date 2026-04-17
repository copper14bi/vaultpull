// Package config loads and validates vaultpull configuration.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddr    string        `yaml:"vault_addr"`
	VaultToken   string        `yaml:"vault_token"`
	SecretPaths  []string      `yaml:"secret_paths"`
	OutputFile   string        `yaml:"output_file"`
	MergeMode    string        `yaml:"merge_mode"`    // overwrite | keep
	RotationInterval string   `yaml:"rotation_interval"` // e.g. "24h", "7d"
	MaxBackups   int           `yaml:"max_backups"`
	AuditLog     string        `yaml:"audit_log"`
	Timeout      time.Duration `yaml:"timeout"`
	Filter       FilterConfig  `yaml:"filter"`
}

// FilterConfig holds key inclusion/exclusion patterns.
type FilterConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// Load reads config from a YAML file, applying defaults.
func Load(path string) (*Config, error) {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}

	return cfg, nil
}

// Validate returns an error if required fields are missing.
func (c *Config) Validate() error {
	if c.VaultToken == "" {
		return errors.New("vault_token is required")
	}
	if len(c.SecretPaths) == 0 {
		return errors.New("at least one secret_path is required")
	}
	return nil
}

func defaults() *Config {
	return &Config{
		VaultAddr:        "http://127.0.0.1:8200",
		OutputFile:       ".env",
		MergeMode:        "overwrite",
		MaxBackups:       5,
		RotationInterval: "24h",
		Timeout:          10 * time.Second,
	}
}
