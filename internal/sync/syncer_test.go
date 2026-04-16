package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

func testConfig(t *testing.T, envFile string, overwrite bool) *config.Config {
	t.Helper()
	return &config.Config{
		VaultAddr: "http://127.0.0.1:8200",
		Token:     "test-token",
		Mappings: []config.Mapping{
			{SecretPath: "secret/app", EnvFile: envFile, Overwrite: overwrite},
		},
	}
}

func TestRun_WritesEnvFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	mock := vault.NewMockClient(map[string]string{"DB_URL": "postgres://localhost/db"})
	cfg := testConfig(t, envFile, true)
	s := sync.NewWithClient(cfg, mock, false)

	results, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Count != 1 {
		t.Errorf("expected count 1, got %d", results[0].Count)
	}
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Error("env file was not created")
	}
}

func TestRun_ErrorOnVaultFailure(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	mock := vault.NewMockClient(nil) // nil triggers error
	cfg := testConfig(t, envFile, true)
	s := sync.NewWithClient(cfg, mock, false)

	_, err := s.Run()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestSummary(t *testing.T) {
	results := []sync.Result{
		{SecretPath: "secret/app", EnvFile: ".env", Count: 3},
	}
	lines := sync.Summary(results)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] == "" {
		t.Error("summary line should not be empty")
	}
}
