package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the application configuration loaded from file or environment.
type Config struct {
	VaultAddr  string            `mapstructure:"vault_addr"`
	VaultToken string            `mapstructure:"vault_token"`
	VaultMount string            `mapstructure:"vault_mount"`
	SecretPaths []string         `mapstructure:"secret_paths"`
	OutputFile string            `mapstructure:"output_file"`
	EnvMapping map[string]string `mapstructure:"env_mapping"`
}

// Load reads configuration from the given file path, falling back to
// environment variables prefixed with VAULTPULL_.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("vault_mount", "secret")
	v.SetDefault("output_file", ".env")

	v.SetEnvPrefix("VAULTPULL")
	v.AutomaticEnv()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultpull")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(os.Getenv("HOME"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	if cfg.VaultToken == "" {
		cfg.VaultToken = os.Getenv("VAULT_TOKEN")
	}

	return &cfg, nil
}

// Validate returns an error if required fields are missing.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault_addr is required")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if len(c.SecretPaths) == 0 {
		return fmt.Errorf("at least one secret_path is required")
	}
	return nil
}
