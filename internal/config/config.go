package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Mapping defines a single Vault path -> env file relationship.
type Mapping struct {
	SecretPath string `yaml:"secret_path"`
	EnvFile    string `yaml:"env_file"`
	Overwrite  bool   `yaml:"overwrite"`
}

// Config holds all vaultpull configuration.
type Config struct {
	VaultAddr string    `yaml:"vault_addr"`
	Token     string    `yaml:"token"`
	Mappings  []Mapping `yaml:"mappings"`
}

// Load reads config from the given file path, applying defaults.
func Load(path string) (*Config, error) {
	cfg := &Config{
		VaultAddr: "http://127.0.0.1:8200",
	}

	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.VaultAddr = addr
	}
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		cfg.Token = token
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// Validate checks that the config has required fields.
func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("vault token is required (set token in config or VAULT_TOKEN env var)")
	}
	if len(c.Mappings) == 0 {
		return errors.New("at least one mapping is required")
	}
	for i, m := range c.Mappings {
		if m.SecretPath == "" {
			return fmt.Errorf("mapping[%d]: secret_path is required", i)
		}
		if m.EnvFile == "" {
			return fmt.Errorf("mapping[%d]: env_file is required", i)
		}
	}
	return nil
}
