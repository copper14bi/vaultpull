package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vaultpull/vaultpull/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	// Ensure no stray env vars interfere.
	t.Setenv("VAULTPULL_VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	cfg, err := config.Load("")
	require.NoError(t, err)

	assert.Equal(t, "http://127.0.0.1:8200", cfg.VaultAddr)
	assert.Equal(t, "secret", cfg.VaultMount)
	assert.Equal(t, ".env", cfg.OutputFile)
	assert.Equal(t, "test-token", cfg.VaultToken)
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "vaultpull.yaml")

	content := `
vault_addr: "http://vault.example.com:8200"
vault_token: "s.abc123"
vault_mount: "kv"
secret_paths:
  - "myapp/prod"
output_file: "prod.env"
`
	err := os.WriteFile(cfgPath, []byte(content), 0600)
	require.NoError(t, err)

	cfg, err := config.Load(cfgPath)
	require.NoError(t, err)

	assert.Equal(t, "http://vault.example.com:8200", cfg.VaultAddr)
	assert.Equal(t, "s.abc123", cfg.VaultToken)
	assert.Equal(t, "kv", cfg.VaultMount)
	assert.Equal(t, []string{"myapp/prod"}, cfg.SecretPaths)
	assert.Equal(t, "prod.env", cfg.OutputFile)
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		SecretPaths: []string{"app/secrets"},
	}
	err := cfg.Validate()
	assert.ErrorContains(t, err, "vault_token is required")
}

func TestValidate_MissingSecretPaths(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "s.token",
	}
	err := cfg.Validate()
	assert.ErrorContains(t, err, "secret_path")
}

func TestValidate_Valid(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		VaultToken:  "s.token",
		SecretPaths: []string{"app/secrets"},
	}
	assert.NoError(t, cfg.Validate())
}
